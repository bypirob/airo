package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bypirob/airo/src/internal/docker"
)

var (
	releaseTag     string
	releaseContext string
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Build, push, and deploy",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		if cfg.Deploy.Type != "ssh" {
			return fmt.Errorf("deploy.type must be ssh for release")
		}
		if releaseTag == "" {
			defaultTag, err := docker.DefaultTag(cfg, projectPath)
			if err != nil {
				return err
			}
			releaseTag = defaultTag
		}

		if err := docker.BuildImage(cfg, projectPath, releaseTag, releaseContext); err != nil {
			return fmt.Errorf("build failed: %w", err)
		}
		if err := docker.PushImage(cfg, projectPath, releaseTag); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}
		if err := docker.Deploy(cfg, releaseTag); err != nil {
			return fmt.Errorf("deploy failed: %w", err)
		}

		return nil
	},
}

func init() {
	releaseCmd.Flags().StringVar(&releaseTag, "tag", "", "image tag (default: <name>:<yyyymmdd-hhmm>-<shortsha>)")
	releaseCmd.Flags().StringVar(&releaseContext, "context", ".", "build context path")
	rootCmd.AddCommand(releaseCmd)
}
