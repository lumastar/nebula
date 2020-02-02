package playlist

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/lumastar/nebula/internal/control"

	"github.com/jleight/omxplayer"
)

type Video struct {
	Path             string  `mapstructure:"path"`
	Loop             bool    `mapstructure:"loop,omitempty"`
	AspectMode       string  `mapstructure:"aspect-mode,omitempty"`
	Orientation      uint    `mapstructure:"orientation,omitempty"`
	Alpha            uint    `mapstructure:"alpha,omitempty"`
	BackgroundColour string  `mapstructure:"background-colour,omitempty"`
	Window           string  `mapstructure:"window,omitempty"`
	Crop             string  `mapstructure:"crop,omitempty"`
	Rate             float32 `mapstructure:"rate,omitempty"`
	// player is the omxplayer instance used by this Video
	player *omxplayer.Player
}

// TODO: Decouple this from the Omxplayer library to allow other implementions for Video to be added in future.

const (
	videoPlayerExecutable = "omxplayer"
)

func (video Video) Check() {
	if video.Path == "" {
		log.Println("Video is missing path", video)
	}
	if !(video.AspectMode == "fill" || video.AspectMode == "stretch" || video.AspectMode == "letterbox") {
		log.Println("Video has invalid aspect mode", video.AspectMode)
	}
	if !(video.Orientation == 0 || video.Orientation == 90 || video.Orientation == 190 || video.Orientation == 270) {
		log.Println("Video has invalid orientation", video.Orientation)
	}
	if video.Alpha > 255 {
		log.Println("Video has invalid alpha", video.Alpha)
	}
	if match, _ := regexp.MatchString("0x[0-9a-f]{8}", video.BackgroundColour); !match {
		log.Println("Video has invalid background colour", video.BackgroundColour)
	}
	// TODO: Check window and crop
	if video.Rate < 0.001 || video.Rate > 4.0 {
		log.Println("Video has invalid rate", video.Rate)
	}
}

func (video Video) Start() {
	// Check the architecture and OS
	if runtime.GOARCH != "arm" || runtime.GOOS != "linux" {
		log.Println("No supported video player for current architecture or OS", runtime.GOARCH, runtime.GOOS)
		return
	}
	// Check that the video player executable can be found
	if _, err := exec.LookPath(videoPlayerExecutable); err != nil {
		log.Println("Unable to find image display executable.", videoPlayerExecutable)
		return
	}
	//omxplayer.SetUser("pi", "/home/pi")
	args := videoPlayerArguments(video)
	player, err := omxplayer.New(video.Path, args...)
	if err != nil {
		log.Println("Error creating player", err)
	}
	video.player = player
	log.Println("Waiting for omxplayer")
	video.player.WaitForReady()
	log.Println("Omxplayer ready, will now play")
	err = player.Play()
	if err != nil {
		log.Println("Error playing video with omxplayer", err)
	}
}

func (video Video) IsRunning() bool {
	if video.player != nil {
		return video.player.IsRunning()
	}
	return false
}

func (video Video) Control(controlItem control.Control) {
	switch controlItem.Command {
	case "stop":
		err := video.player.Stop()
		if err != nil {
			log.Println("Error stopping omxplayer video", err)
		}
	}
}

func videoPlayerArguments(video Video) []string {
	arguments := []string{"--no-osd"}
	if video.Loop == true {
		arguments = append(arguments, "--loop")
	}
	if video.AspectMode != "" {
		arguments = append(arguments, "--aspect-mode")
		arguments = append(arguments, video.AspectMode)
	}
	if video.Orientation != 0 {
		arguments = append(arguments, "--orientation")
		arguments = append(arguments, fmt.Sprintf("%v", video.Orientation))
	}
	if video.Alpha != 0 {
		arguments = append(arguments, "--alpha")
		arguments = append(arguments, fmt.Sprintf("%v", video.Alpha))
	}
	if video.BackgroundColour != "" {
		arguments = append(arguments, "--blank")
		arguments = append(arguments, video.BackgroundColour)
	}
	if video.Window != "" {
		arguments = append(arguments, "--win")
		arguments = append(arguments, video.Window)
	}
	if video.Crop != "" {
		arguments = append(arguments, "--crop")
		arguments = append(arguments, video.Crop)
	}
	if video.Rate != 0 {
		arguments = append(arguments, "--rate")
		arguments = append(arguments, fmt.Sprintf("%v", video.Rate))
	}
	arguments = append(arguments, "--layer")
	//arguments = append(arguments, layer) TODO: Get layer to here
	// Don't add the video path here as that is required separately by the Omxplayer library
	return arguments
}
