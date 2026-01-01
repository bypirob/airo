package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"bypirob/airo/src/internal/docker"
)

var deployTag string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a Docker image over SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		if cfg.Deploy.Type != "ssh" {
			return fmt.Errorf("deploy.type must be ssh for deploy")
		}
		if deployTag == "" {
			return fmt.Errorf("--tag is required")
		}

		if err := docker.Deploy(cfg, deployTag); err != nil {
			return fmt.Errorf("deploy failed: %w", err)
		}

		return nil
	},
}

func init() {
	deployCmd.Flags().StringVar(&deployTag, "tag", "", "image tag to deploy")
	_ = deployCmd.MarkFlagRequired("tag")
	rootCmd.AddCommand(deployCmd)
}
