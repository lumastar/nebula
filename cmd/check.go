package cmd

import (
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check Nebula configuration",
	Long: `Check Nebula configuration to ensure it is valid.
`,
	Run: func(cmd *cobra.Command, args []string) {
		loadConfigFile()
		check()
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.PersistentFlags().StringVarP(
		&configFilePath,
		"config-file",
		"c",
		"",
		"Config file location, without this flag we search for `nebula.yaml` in the current working directory",
	)
}

func check() {
	for _, playlist := range c.Playlists {
		playlist.Check()
	}
}
