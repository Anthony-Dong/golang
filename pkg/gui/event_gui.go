package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
)

const SideView = "side"
const MainView = "main"
const HelpView = "help"

type keyManager struct {
	lastTime time.Time
	hook     func(g *gocui.Gui, v *gocui.View) error
}

func NewKeyManager(hook func(g *gocui.Gui, v *gocui.View) error) *keyManager {
	return &keyManager{
		hook: hook,
	}
}

func (k *keyManager) NewHandler() func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		currentTime := time.Now()
		if currentTime.Sub(k.lastTime) < 1*time.Millisecond {
			return nil
		}
		k.lastTime = currentTime
		return k.hook(g, v)
	}
}

type Event interface {
	GetBasicInfo() string
	GetDetailInfo() string
}

type SimpleEvent struct {
	Basic  string
	Detail string
}

func (s *SimpleEvent) GetBasicInfo() string {
	return s.Basic
}

func (s *SimpleEvent) GetDetailInfo() string {
	return s.Detail
}

type EventGUI struct {
	events []Event
	recv   chan Event

	sideWidth      int
	sideSmallWidth int
	helpViewHeight int

	helpInfo string
	helpLogs []string

	g *gocui.Gui
}

func (e *EventGUI) refreshHelpView() {
	view, err := e.GetGui().View(HelpView)
	if err != nil {
		return
	}
	view.Clear()
	Println(view, e.helpInfo)
	for _, elem := range e.helpLogs {
		Println(view, elem)
	}
}

func (e *EventGUI) setHelpInfo(info string) {
	e.helpInfo = info
	e.refreshHelpView()
}

func (e *EventGUI) log(format string, args ...interface{}) {
	if len(e.helpLogs) == e.helpViewHeight-1 {
		e.helpLogs = e.helpLogs[1:]
	}
	e.helpLogs = append(e.helpLogs, fmt.Sprintf(format, args...))
	e.refreshHelpView()
}

func (e *EventGUI) GetGui() *gocui.Gui {
	return e.g
}

func NewEventGUI() (*EventGUI, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}
	g.Cursor = true
	g.Highlight = true
	x, _ := g.Size()
	r := &EventGUI{
		events:         make([]Event, 0, 1024),
		recv:           make(chan Event, 1024),
		g:              g,
		sideWidth:      int((float64(x) / float64(100)) * float64(20)),
		sideSmallWidth: int((float64(x) / float64(100)) * float64(10)),
		helpViewHeight: 5,
	}
	if err := r.init(); err != nil {
		return nil, err
	}
	return r, nil
}

func (e *EventGUI) Close() {
	e.GetGui().Close()
}

func (e *EventGUI) init() error {
	if err := e.createMainView(nil, true); err != nil {
		panic(err)
	}
	if err := e.createSideView(false); err != nil {
		panic(err)
	}

	if err := e.bindKey(); err != nil {
		panic(err)
	}
	return nil
}

func (e *EventGUI) MainLoop() error {
	go func() {
		for {
			_ = e.renderEvent(<-e.recv)
		}
	}()
	g := e.GetGui()
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

func (e *EventGUI) AddEvent(event Event) {
	e.recv <- event
}

func (e *EventGUI) renderEvent(event Event) error {
	e.GetGui().Update(func(g *gocui.Gui) error {
		defer func() {
			if r := recover(); r != nil {
				e.log("render event panic: %v", r)
			}
		}()
		view, err := g.View(SideView)
		if err != nil {
			if err == gocui.ErrUnknownView {
				return nil
			}
			return err
		}
		e.events = append(e.events, event)
		info := event.GetBasicInfo()
		view.Write([]byte(info))
		view.Write([]byte{'\n'})

		e.cursorDown(g, view, event)

		return nil
	})
	return nil
}

func (e *EventGUI) currentViewIs(view string) bool {
	v := e.GetGui().CurrentView()
	return v.Name() == view
}

func (e *EventGUI) cursorDown(_ *gocui.Gui, v *gocui.View, event Event) {
	if !e.currentViewIs(SideView) {
		return
	}
	if len(e.events) > 0 {
		index := e.getCurrentIndex()
		if index+1 != len(e.events)-1 {
			return
		}
	}
	if lines := len(e.events); lines > 0 {
		if err := v.SetCursor(0, lines-1); err != nil {
			_, oy := v.Origin()
			_ = v.SetOrigin(0, oy+1)
		}
	}
	_ = e.showDetail(event, true)
	_, cy := v.Cursor()
	_, oy := v.Origin()
	e.log("cursor down. size: %d, index: %d, cursor: %d, origin: %d", len(e.events), e.getCurrentIndex(), cy, oy)
}

func (e *EventGUI) showDetail(event Event, simple bool) error {
	if event == nil {
		return nil
	}
	view, err := e.GetGui().View(MainView)
	if err != nil {
		if err == gocui.ErrUnknownView {
			return nil
		}
		return err
	}
	view.Clear()
	if simple {
		_, y := view.Size()
		result := strings.SplitN(event.GetDetailInfo(), "\n", y+1)
		if len(result) == y+1 {
			result = result[:y]
		}
		Printf(view, strings.Join(result, "\n"))
	} else {
		Printf(view, event.GetDetailInfo())
	}
	view.Write([]byte{'\n'})

	return nil
}

func (e *EventGUI) bindKey() error {
	g := e.GetGui()
	if err := g.SetKeybinding(SideView, gocui.KeyArrowDown, gocui.ModNone, NewKeyManager(e.nextEvent).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding(SideView, gocui.KeyArrowUp, gocui.ModNone, NewKeyManager(e.prevEvent).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding(MainView, gocui.KeyArrowLeft, gocui.ModNone, NewKeyManager(func(g *gocui.Gui, v *gocui.View) error {
		vx, vy := v.Cursor()
		ox, oy := v.Origin()
		if vx+ox == 0 {
			return nil
		}
		if err := v.SetCursor(vx-1, vy); err != nil {
			return v.SetOrigin(ox-1, oy)
		}
		return nil
	}).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding(MainView, gocui.KeyArrowRight, gocui.ModNone, NewKeyManager(func(g *gocui.Gui, v *gocui.View) error {
		vx, vy := v.Cursor()
		ox, oy := v.Origin()
		if err := v.SetCursor(vx+1, vy); err != nil {
			return v.SetOrigin(ox+1, oy)
		}
		return nil
	}).NewHandler()); err != nil {
		return err
	}

	if err := g.SetKeybinding(MainView, gocui.KeyArrowUp, gocui.ModNone, NewKeyManager(func(g *gocui.Gui, v *gocui.View) error {
		vx, vy := v.Cursor()
		ox, oy := v.Origin()
		if vy+oy == 0 {
			return nil
		}
		if err := v.SetCursor(vx, vy-1); err != nil {
			return v.SetOrigin(ox, oy-1)
		}
		return nil
	}).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding(MainView, gocui.KeyArrowDown, gocui.ModNone, NewKeyManager(func(g *gocui.Gui, v *gocui.View) error {
		vx, vy := v.Cursor()
		ox, oy := v.Origin()
		if err := v.SetCursor(vx, vy+1); err != nil {
			return v.SetOrigin(ox, oy+1)
		}
		return nil
	}).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding(MainView, gocui.KeyF1, gocui.ModNone, NewKeyManager(func(g *gocui.Gui, v *gocui.View) error {
		if err := clipboard.WriteAll(v.Buffer()); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		e.log("拷贝成功!")
		return nil
	}).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, NewKeyManager(func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}).NewHandler()); err != nil {
		return err
	}

	if err := g.SetKeybinding(SideView, gocui.KeyEnter, gocui.ModNone, NewKeyManager(e.switchView).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding(SideView, gocui.KeyArrowRight, gocui.ModNone, NewKeyManager(e.switchView).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding(SideView, gocui.KeyTab, gocui.ModNone, NewKeyManager(e.switchView).NewHandler()); err != nil {
		return err
	}
	if err := g.SetKeybinding(MainView, gocui.KeyTab, gocui.ModNone, NewKeyManager(e.switchView).NewHandler()); err != nil {
		return err
	}
	return nil
}

func (e *EventGUI) createSideView(small bool) error {
	g := e.GetGui()
	_, maxY := g.Size()
	x1 := e.sideWidth
	if small {
		x1 = e.sideSmallWidth
	}
	if view, err := g.SetView(SideView, -1, -1, x1, maxY-e.helpViewHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "side"
		view.Highlight = true
		//view.Autoscroll = true
		view.SelBgColor = gocui.ColorGreen
		view.SelFgColor = gocui.ColorBlack
	}
	e.GetGui().SetCurrentView(SideView)
	e.helpFlag("[列表页] Tab: 切换到详情页")
	if event := e.getCurrentEvent(); event != nil {
		e.showDetail(event, true)
	}
	return nil
}

func (e *EventGUI) createMainView(event Event, small bool) error {
	g := e.GetGui()
	maxX, maxY := g.Size()
	x0 := e.sideSmallWidth
	if small {
		x0 = e.sideWidth
	}
	if view, err := g.SetView(MainView, x0, -1, maxX, maxY-e.helpViewHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "main"
	}
	g.SetCurrentView(MainView)
	e.showDetail(event, false)
	e.helpFlag("[详情页] F1: 拷贝详情\tTab: 切换到列表页")
	return nil
}

func (e *EventGUI) switchView(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == MainView {
		g.DeleteView(MainView)
		e.createMainView(nil, true)
		return e.createSideView(false)
	}
	if v.Name() == SideView {
		event := e.getCurrentEvent()
		if event == nil {
			return nil
		}
		e.createSideView(true)
		return e.createMainView(event, false)
	}
	return nil
}

func (e *EventGUI) bindCopyDashboard(g *gocui.Gui) error {
	mm := NewKeyManager(func(g *gocui.Gui, v *gocui.View) error {
		if err := clipboard.WriteAll(v.Buffer()); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		return nil
	})
	return g.SetKeybinding(MainView, gocui.KeyF1, gocui.ModNone, mm.NewHandler())
}

func (e *EventGUI) prevEvent(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	if cy == 0 && oy == 0 {
		return nil
	}
	if err := v.SetCursor(cx, cy-1); err != nil {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	event := e.getCurrentEvent()
	if event == nil {
		return nil
	}
	return e.showDetail(event, true)
}

func (e *EventGUI) nextEvent(g *gocui.Gui, v *gocui.View) error {
	if index := e.getCurrentIndex(); index == len(e.events)-1 {
		return nil
	}
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	if err := v.SetCursor(cx, cy+1); err != nil {
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	event := e.getCurrentEvent()
	if event == nil {
		return nil
	}
	return e.showDetail(event, true)
}

func (e *EventGUI) getCurrentEvent() Event {
	index := e.getCurrentIndex()
	if index >= len(e.events) {
		return nil
	}
	return e.events[index]
}

func (e *EventGUI) getCurrentIndex() int {
	v, err := e.GetGui().View(SideView)
	if err != nil {
		return -1
	}
	_, cy := v.Cursor() // cursor index
	_, oy := v.Origin() // buffer index
	if oy == 0 {
		return cy
	}
	return cy + oy
}

type HelpInfo struct {
	ViewName string
	Flags    map[string]string
}

func (e *EventGUI) helpFlag(help string) error {
	g := e.GetGui()
	maxX, maxY := g.Size()
	if view, err := g.SetView(HelpView, -1, maxY-e.helpViewHeight, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Title = "help"
		e.setHelpInfo(help)
		e.refreshHelpView()
	} else {
		e.setHelpInfo(help)
		e.refreshHelpView()
	}
	return nil
}
