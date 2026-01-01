package docker

import (
	"fmt"
	"strings"

	"bypirob/airo/src/internal/config"
)

func Status(cfg config.Config) (string, error) {
	running, err := dockerStatus(cfg, false)
	if err != nil {
		return "", err
	}
	if running != "" {
		return "running", nil
	}

	any, err := dockerStatus(cfg, true)
	if err != nil {
		return "", err
	}
	if any != "" {
		return "stopped", nil
	}

	return "not found", nil
}

func dockerStatus(cfg config.Config, all bool) (string, error) {
	args := []string{"docker", "ps"}
	if all {
		args = append(args, "-a")
	}
	args = append(args, "--filter", fmt.Sprintf("name=^%s$", cfg.Name), "--format", "{{.Status}}")

	cmd := sshCommand(cfg, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ssh status: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
