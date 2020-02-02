package control

type Control struct {
	Command string   `mapstructure:"command"`
	Value   string   `mapstructure:"value,omitempty"`
	Targets []Target `mapstructure:"targets,omitempty"`
}

type Target struct {
	Type   string `mapstructure:"type"`
	Name   string `mapstructure:"name,omitempty"`
	Number uint   `mapstructure:"number,omitempty"`
}

var (
	commands = [...]string{"skip", "stop", "quit", "goto", "alpha", "on", "off", "blank", "show", "next", "rate"}
)

func (control Control) Check() {
	// TODO: Check whether command is recognised.
	switch control.Command {
	case "skip":
	case "stop":
	case "alpha":
		// TODO: Check value is given, convert to unint, check between 0 and 255.
	default:
		// ERROR
	}
	// TODO: Check whether commands that require values have a value.
	// TODO: Check whether given value is valid for command.
	// TODO: Check wheth
}

func (control Control) Start() {
	// TODO: Issue the command
}

func (control Control) IsRunning() bool {
	return false
}

func (control Control) Control(controlItem Control) {
}
