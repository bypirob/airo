package main

import (
	"github.com/spf13/cobra"

	"bypirob/airo/src/internal/docker"
)

var tagsRemote bool

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List image tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		tags, err := docker.Tags(cfg, tagsRemote)
		if err != nil {
			return err
		}

		for _, tag := range tags {
			cmd.Println(tag)
		}

		return nil
	},
}

func init() {
	tagsCmd.Flags().BoolVar(&tagsRemote, "remote", false, "list tags from the registry")
	rootCmd.AddCommand(tagsCmd)
}
