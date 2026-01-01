package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bypirob/airo/src/internal/docker"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check remote container status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		if cfg.Deploy.Type != "ssh" {
			return fmt.Errorf("deploy.type must be ssh for status")
		}

		status, err := docker.Status(cfg)
		if err != nil {
			return err
		}

		cmd.Println(status)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
