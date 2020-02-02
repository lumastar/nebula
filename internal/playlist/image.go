package playlist

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"

	"github.com/lumastar/nebula/internal/control"
)

type Image struct {
	Path             string  `mapstructure:"path"`
	Offset           Offset  `mapstructure:"offset,omitempty"`
	Sprite           Sprite  `mapstructure:"sprite,omitempty"`
	BackgroundColour string  `mapstructure:"background-colour,omitempty"`
	Duration         float32 `mapstructure:"duration,omitempty"`
	// command is the Spriteview command instance used by this Image
	command *exec.Cmd
}

type Offset struct {
	X int `mapstructure:"x,omitempty"`
	Y int `mapstructure:"y,omitempty"`
}

type Sprite struct {
	Rows     uint    `mapstructure:"rows,omitempty"`
	Columns  uint    `mapstructure:"columns,omitempty"`
	Interval float32 `mapstructure:"interval,omitempty"`
}

const (
	imageDisplayExecutable = "spriteview"
)

func (image Image) Check() {
	if image.Path == "" {
		log.Println("Video is missing path", image)
	}
	if image.Sprite.Interval > 0.0 {
		log.Println("Video has invalid sprite interval", image.Sprite.Interval)
	}
	if match, _ := regexp.MatchString("0x[0-9a-f]{8}", image.BackgroundColour); !match {
		log.Println("Video has invalid background colour", image.BackgroundColour)
	}
	if image.Duration < 0.0 {
		log.Println("Video has invalid duration", image.Duration)
	}
}

func (image Image) Start() {
	if _, err := exec.LookPath(imageDisplayExecutable); err != nil {
		log.Println("Unable to find image display executable.", imageDisplayExecutable)
		return
	}
	args := imageDisplayArguments(image)
	command := exec.Command(imageDisplayExecutable, args...)
	err := command.Run()
	if err != nil {
		log.Println("Failed to run image display executable.", imageDisplayExecutable, err)
	}
	image.command = command
}

func (image Image) IsRunning() bool {
	return false
}

func (image Image) Control(control control.Control) {
}

func imageDisplayArguments(image Image) []string {
	arguments := []string{}
	if image.Offset.X != 0 {
		arguments = append(arguments, "-x")
		arguments = append(arguments, fmt.Sprintf("%v", image.Offset.X))
	}
	if image.Offset.Y != 0 {
		arguments = append(arguments, "-y")
		arguments = append(arguments, fmt.Sprintf("%v", image.Offset.Y))
	}
	if image.Sprite.Rows != 0 {
		arguments = append(arguments, "-r")
		arguments = append(arguments, fmt.Sprintf("%v", image.Sprite.Rows))
	}
	if image.Sprite.Columns != 0 {
		arguments = append(arguments, "-c")
		arguments = append(arguments, fmt.Sprintf("%v", image.Sprite.Columns))
	}
	if image.Sprite.Interval != 0 {
		arguments = append(arguments, "-i")
		arguments = append(arguments, fmt.Sprintf("%v", image.Sprite.Interval))
	}
	// Unless sprite rows or columns are set with no duration then run in non interactive mode
	if !((image.Sprite.Rows != 0 || image.Sprite.Columns != 0) && image.Sprite.Interval == 0) {
		arguments = append(arguments, "-n")
	}
	if image.Duration != 0 {
		arguments = append(arguments, "-t")
		arguments = append(arguments, fmt.Sprintf("%v", image.Duration))
	}
	arguments = append(arguments, "-l")
	//arguments = append(arguments, layer) TODO: Get layer to here
	arguments = append(arguments, image.Path)
	return arguments
}
