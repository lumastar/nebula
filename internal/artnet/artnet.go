package artnet

const (
	MAX_UNIVERSE = 10
	MAX_CHANNEL  = 255
)

type ArtNetUniverseChannel struct {
	Universe uint
	Channel  uint
	Inputs   []interface{}
}

func Check(universeChannels []ArtNetUniverseChannel) {
	for _, universeChannel := range universeChannels {
		if universeChannel.Universe > MAX_UNIVERSE {

		}
		if universeChannel.Channel > MAX_CHANNEL {

		}
	}
}

// Listen for Art-Net input and handle
func Listen() {
}
