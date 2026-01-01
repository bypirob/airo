package main

import (
	"github.com/spf13/cobra"

	"bypirob/airo/src/internal/config"
)

var (
	projectPath string
	configPath  string
)

var rootCmd = &cobra.Command{
	Use:          "airo",
	Short:        "airo builds and deploys container images",
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&projectPath, "project", ".", "path to project directory")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "airo.yaml", "config file path, relative to --project")
}

func loadConfig() (config.Config, error) {
	return config.Load(projectPath, configPath)
}
