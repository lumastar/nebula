package playlist

import "github.com/lumastar/nebula/internal/control"

type Label struct {
	Name string `mapstructure:"name"`
	Sync bool   `mapstructure:"sync,omitempty"`
}

func (label Label) Check() {
}

func (label Label) Start() {
}

func (label Label) IsRunning() bool {
	return false
}

func (label Label) Control(control control.Control) {
}
