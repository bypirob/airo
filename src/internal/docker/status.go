package docker

import (
	"fmt"
	"strings"

	"bypirob/airo/src/internal/config"
)

func Status(cfg config.Config) (string, error) {
	anyRunning := false
	anyFound := false
	for _, container := range cfg.Deploy.Containers {
		running, err := dockerStatus(cfg, container.Name, false)
		if err != nil {
			return "", err
		}
		if running != "" {
			anyRunning = true
			anyFound = true
			continue
		}

		any, err := dockerStatus(cfg, container.Name, true)
		if err != nil {
			return "", err
		}
		if any != "" {
			anyFound = true
		}
	}

	if anyRunning {
		return "running", nil
	}
	if anyFound {
		return "stopped", nil
	}
	return "not found", nil
}

func dockerStatus(cfg config.Config, name string, all bool) (string, error) {
	args := []string{"docker", "ps"}
	if all {
		args = append(args, "-a")
	}
	args = append(args, "--filter", fmt.Sprintf("name=^%s$", name), "--format", "{{.Status}}")

	cmd := sshCommand(cfg, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ssh status: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
