package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/lumastar/nebula/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "nebula",
	Short: "Media playback controller",
	Long: `Nebula is an application for controlling video playback and image display,
designed for use with projection mapping, visual shows, and large format displays.
Layered media playlists and playback controls are specified in a config file.
Nebula can then be operated from a command prompt, Art-Net, and GPIO.
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var configFilePath string
var c config.Config

// TODO: All the load and check functionality should move to a separate cmd and be called by this cmd
func loadConfigFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigFile(configFilePath)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		// Not having a configuration file is an usual case, so alert on it.
		log.Println("Error reading config file:", viper.ConfigFileUsed())
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Println("Unable to decode into struct", err)
	}
}
