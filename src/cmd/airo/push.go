package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bypirob/airo/src/internal/docker"
)

var pushTag string

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push a Docker image",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		if err := docker.PushImage(cfg, projectPath, pushTag); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}

		return nil
	},
}

func init() {
	pushCmd.Flags().StringVar(&pushTag, "tag", "", "image tag (default: <name>:<yyyymmdd-hhmm>-<shortsha>)")
	rootCmd.AddCommand(pushCmd)
}
