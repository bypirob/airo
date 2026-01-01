package main

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(versionString())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionString() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	if version == "" {
		return "dev"
	}
	return version
}
