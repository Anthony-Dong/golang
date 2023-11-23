package tools

type Dlv struct {
	ProgramArgs []string
	Action      string
	DlvArgs     []string
}

func NewDebugDlv(pkg string, args ...string) *Dlv {
	return &Dlv{
		ProgramArgs: args,
		Action:      "debug",
		DlvArgs:     []string{pkg},
	}
}

func NewTestDlv(testName string, pkg string) *Dlv {
	return &Dlv{
		ProgramArgs: []string{
			"-test.v",
			"-test.run",
			testName,
			"-test.count",
			"1",
		},
		Action:  "test",
		DlvArgs: []string{pkg},
	}
}

func (d *Dlv) Name() string {
	return "dlv"
}

func (d *Dlv) Args() []string {
	args := make([]string, 0)
	args = append(args, d.Action)
	//args = append(args, `--log`)
	args = append(args, `--headless`)
	args = append(args, `--listen=:10086`)
	args = append(args, `--api-version=2`)
	args = append(args, `--accept-multiclient`)
	args = append(args, `--build-flags=-v -gcflags 'all=-N -l'`)
	for _, elem := range d.DlvArgs {
		args = append(args, elem)
	}
	if len(d.ProgramArgs) > 0 {
		args = append(args, "--")
	}
	for _, elem := range d.ProgramArgs {
		args = append(args, elem)
	}
	return args
}
