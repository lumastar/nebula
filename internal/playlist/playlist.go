package playlist

import (
	"log"
	"time"

	"github.com/lumastar/nebula/internal/control"

	"github.com/mitchellh/mapstructure"
)

// Playlist is a list of media items for a specified layer.
type Playlist struct {
	Name  string        `mapstructure:"name"`
	Layer uint          `mapstructure:"layer,omitempty"`
	Items []interface{} `mapstructure:"items,omitempty"`
	// controlInput is the channel used to send control items to a running Playlist.
	// The Control function adds items to the channel, which is then read as part of the Run loop.
	controlInput chan control.Control
	// controlOutput
	controlOutput chan control.Control
}

// TODO: It might be better to make a separate PlaylistRun struct to hold a Playlist with controlInput and controlOutput.

// Item is a media item in the Playlist.
type Item interface {
	// Check determines whether the items values are valid.
	// E.g. for a video is the path set, is the aspect mode a recognised value, etc.
	Check()
	// Start should start the item and return.
	// E.g. for a video this will start the playback.
	// Start may launch a run function in a thread if the item required continuous monitoring.
	Start()
	// Is running should return true if the item is still running.
	// E.g. for a video this will return true if it is still playing.
	IsRunning() bool
	// Control the item with a control item.
	// E.g. for a video if the control item has command 'pause' then pause playback.
	Control(controlItem control.Control)
}

// Check validates the Playlist and all its Items.
func (playlist Playlist) Check() {
	for _, encodedItem := range playlist.Items {
		item := decodeItem(encodedItem)
		item.Check()
	}
	// TODO: Check for duplicate label names in Playlist
}

// Start sets up the Playlist and calls its Run function in another thread.
func (playlist Playlist) Start(controlOutput chan control.Control) {
	playlist.controlOutput = controlOutput
	playlist.controlInput = make(chan control.Control, 10)
	go playlist.run()
}

// run is the function to start the Playlist, designed to be run as a separate thread.
func (playlist Playlist) run() {
	for _, encodedItem := range playlist.Items {
		item := decodeItem(encodedItem)
		// TODO: How will Label and Control Items work?
		// They need access to the controlOutput, but the other types do not.
		item.Start()
		for item.IsRunning() {
			select {
			case controlItem := <-playlist.controlInput:
				switch controlItem.Command {
				case "skip":
				case "stop":
				// Only playlist level commands need to be handled here.
				// Otherwise they can just be passed on to the item.
				// E.g. the 'alpha' command can just be handled by the current Video or Image.
				default:
					item.Control(controlItem)
				}
			default:
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// Control is a function to send control items to a running Playlist.
// It adds them to a channel that is read by the Run function, running in a separate thread.
func (playlist Playlist) Control(controlItem control.Control) {
	// TODO: This will probably go badly if called before Start
	playlist.controlInput <- controlItem
}

// decodeItem determines the type of a Playlist item specified in the configuration file.
// It then decodes it with the correct struct and returns it.
func decodeItem(encodedItem interface{}) Item {
	type TypeItem struct {
		Type string
	}
	var typeItem TypeItem
	err := mapstructure.Decode(encodedItem, &typeItem)
	if err != nil {
		log.Println("Unable to determine type of item", err)
	}
	var item Item
	switch typeItem.Type {
	case "label":
		var label Label
		err = mapstructure.Decode(item, &label)
		item = label
	case "video":
		var video Video
		err = mapstructure.Decode(item, &video)
		item = video
	case "image":
		var image Image
		err = mapstructure.Decode(item, &image)
		item = image
	case "control":
		var control control.Control
		err = mapstructure.Decode(item, &control)
		item = control
	default:
		log.Println("Item has unknown type", typeItem.Type)
	}
	if err != nil {
		log.Println("Unable to decode into struct", err)
	}
	return item
}
