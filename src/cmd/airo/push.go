package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bypirob/airo/src/internal/docker"
)

var pushCmd = &cobra.Command{
	Use:   "push <tag>",
	Short: "Push a Docker image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		tag := args[0]
		if err := docker.PushImage(cfg, projectPath, tag); err != nil {
			return fmt.Errorf("push failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
