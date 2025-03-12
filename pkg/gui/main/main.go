package main

import (
	"fmt"
	"time"

	"github.com/anthony-dong/golang/pkg/gui"
)

func main() {
	g, err := gui.NewEventGUI()
	if err != nil {
		panic(err)
	}
	defer g.Close()
	go func() {
		for x := 0; x < 1000000; x++ {
			g.AddEvent(&gui.SimpleEvent{
				Basic:  fmt.Sprintf("item %d basic info 1111  111111    1111", x),
				Detail: fmt.Sprintf("你好、\nitem %d detail info\n 222333 \n 44444", x),
			})
			time.Sleep(time.Millisecond * 100)
		}
	}()
	if err := g.MainLoop(); err != nil {
		panic(err)
	}
}
