package docker

import (
	"fmt"
	"os"

	"bypirob/airo/src/internal/config"
)

func Deploy(cfg config.Config, tag string) error {
	tags, err := resolveTags(cfg, "", tag)
	if err != nil {
		return err
	}

	for _, container := range cfg.Deploy.Containers {
		imageTag := tags[container.Image]
		runArgs := []string{"docker", "run", "-d", "--name", container.Name}
		if container.Port != 0 && container.AppPort != 0 {
			runArgs = append(runArgs, "-p", fmt.Sprintf("%d:%d", container.Port, container.AppPort))
		}
		if container.EnvFile != "" {
			runArgs = append(runArgs, "--env-file", container.EnvFile)
		}
		for _, network := range container.Networks {
			if network == "" {
				continue
			}
			runArgs = append(runArgs, "--network", network)
		}
		runArgs = append(runArgs, imageTag)

		stopCmd := shellJoin([]string{"docker", "stop", container.Name})
		removeCmd := shellJoin([]string{"docker", "rm", "-f", container.Name})
		runCmd := shellJoin(runArgs)
		remoteCmd := fmt.Sprintf("%s >/dev/null 2>&1 || true; %s >/dev/null 2>&1 || true; %s", stopCmd, removeCmd, runCmd)

		sshCmd := sshCommand(cfg, "sh", "-c", remoteCmd)
		sshCmd.Stdout = os.Stdout
		sshCmd.Stderr = os.Stderr

		if err := sshCmd.Run(); err != nil {
			return fmt.Errorf("ssh deploy (%s): %w", container.Name, err)
		}
	}

	return nil
}
