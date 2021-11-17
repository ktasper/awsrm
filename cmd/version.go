package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version information",
	Long:  `Version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func version() string {
	return "1.1.0"
}
