package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bypirob/airo/src/internal/docker"
)

var (
	buildTag     string
	buildContext string
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a Docker image",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		if err := docker.BuildImage(cfg, projectPath, buildTag, buildContext); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}

		return nil
	},
}

func init() {
	buildCmd.Flags().StringVar(&buildTag, "tag", "", "image tag suffix (default: <yyyymmdd-hhmm>-<shortsha>)")
	buildCmd.Flags().StringVar(&buildContext, "context", ".", "build context path")
	rootCmd.AddCommand(buildCmd)
}
