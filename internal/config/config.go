package config

import (
	"github.com/lumastar/nebula/internal/playlist"
)

// Config is the full configuration for Nebula to run
type Config struct {
	Playlists []playlist.Playlist
	//Projectors []Projector
	//ArtNet     []ArtNet
	//GPIO       []GPIO
}

// Projector is a projector to be controlled over PJLink
type Projector struct {
	IP       string
	Password string
	Name     string
	Number   uint
}

type GPIO struct {
}
