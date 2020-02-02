package cmd

import (
	"github.com/lumastar/nebula/internal/control"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Nebula",
	Long: `Run Nebula.
`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfigFile()
		run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringVarP(
		&configFilePath,
		"config-file",
		"c",
		"",
		"Config file location, without this flag we search for `nebula.yaml` in the current working directory",
	)
}

func run() {
	// controlOutput is used by Playlists to output Control items.
	controlOutput := make(chan control.Control, 10)
	// Start each of the Playlists.
	for _, playlist := range c.Playlists {
		playlist.Start(controlOutput)
	}
	// TODO: Start each of the Projectors.
	// TODO: Start Art-Net.
	// TODO: Start GPIO.
	for {
		controlItem := <-controlOutput
		// If no Targets are specified send this controlItem to all Playlists, Projectors, etc.
		if len(controlItem.Targets) == 0 {
			for _, playlist := range c.Playlists {
				playlist.Control(controlItem)
			}
			// TODO: Send to Projectors
		} else {
			for _, target := range controlItem.Targets {
				if target.Type == "Playlist" {
					for _, playlist := range c.Playlists {
						// If the Target name is not set send to all Playlists.
						// Otherwise only send if the Target Name matches the Playlist Name.
						if target.Name == "" || target.Name == playlist.Name {
							playlist.Control(controlItem)
						}
					}
				} else if target.Type == "Projector" {
					// TODO: Send to Projectors
				}
			}
		}
	}
}
