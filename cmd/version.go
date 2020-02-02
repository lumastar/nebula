package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var verbose bool

// These variables are injected at build time.

// nebulaVersion is the version of the app.
var nebulaVersion = "development"

// commit is the commit hash of the build
var commit string

// buildDate is the date it was built
var buildDate string

// goVersion is the go version that was used to compile this
var goVersion string

// platform is the target platform this was compiled for
var platform string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version",
	Long: `Display Nebula version.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nebula version: ", nebulaVersion, platform)
		if verbose {
			fmt.Println()
			fmt.Println("Commit: ", commit)
			fmt.Println("Built: ", buildDate)
			fmt.Println("Go: ", goVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"If enabled, displays the additional information about this built.",
	)
}
