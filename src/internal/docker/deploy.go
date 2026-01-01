package docker

import (
	"fmt"
	"os"

	"bypirob/airo/src/internal/config"
)

func Deploy(cfg config.Config, tag string) error {
	if tag == "" {
		return fmt.Errorf("tag is required for deploy")
	}

	runArgs := []string{"docker", "run", "-d", "--name", cfg.Name}
	if cfg.Container.Port != 0 && cfg.Container.AppPort != 0 {
		runArgs = append(runArgs, "-p", fmt.Sprintf("%d:%d", cfg.Container.Port, cfg.Container.AppPort))
	}
	if cfg.Deploy.EnvFile != "" {
		runArgs = append(runArgs, "--env-file", cfg.Deploy.EnvFile)
	}
	for _, network := range cfg.Deploy.Networks {
		if network == "" {
			continue
		}
		runArgs = append(runArgs, "--network", network)
	}
	runArgs = append(runArgs, tag)

	stopCmd := shellJoin([]string{"docker", "stop", cfg.Name})
	removeCmd := shellJoin([]string{"docker", "rm", "-f", cfg.Name})
	runCmd := shellJoin(runArgs)
	remoteCmd := fmt.Sprintf("%s >/dev/null 2>&1 || true; %s >/dev/null 2>&1 || true; %s", stopCmd, removeCmd, runCmd)

	sshCmd := sshCommand(cfg, "sh", "-c", remoteCmd)
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Run(); err != nil {
		return fmt.Errorf("ssh deploy: %w", err)
	}

	return nil
}
