package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bypirob/airo/src/internal/config"
)

func PushImage(cfg config.Config, projectPath, tag string) error {
	if projectPath == "" {
		projectPath = "."
	}

	tags, err := resolveTags(cfg, projectPath, tag)
	if err != nil {
		return err
	}

	switch cfg.Deploy.Type {
	case "ssh":
		return pushOverSSH(cfg, tags)
	case "registry":
		return pushToRegistry(cfg, tags)
	default:
		return fmt.Errorf("unsupported deploy.type %q", cfg.Deploy.Type)
	}
}

func pushOverSSH(cfg config.Config, tags map[string]string) error {
	for name, tag := range tags {
		saveCmd := exec.Command("docker", "save", tag)
		sshCmd := sshCommand(cfg, "docker", "load")

		pipe, err := saveCmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("prepare docker save output (%s): %w", name, err)
		}

		saveCmd.Stderr = os.Stderr
		sshCmd.Stdin = pipe
		sshCmd.Stdout = os.Stdout
		sshCmd.Stderr = os.Stderr

		if err := sshCmd.Start(); err != nil {
			return fmt.Errorf("start ssh (%s): %w", name, err)
		}
		if err := saveCmd.Start(); err != nil {
			return fmt.Errorf("start docker save (%s): %w", name, err)
		}

		saveErr := saveCmd.Wait()
		sshErr := sshCmd.Wait()

		if saveErr != nil {
			return fmt.Errorf("docker save (%s): %w", name, saveErr)
		}
		if sshErr != nil {
			return fmt.Errorf("ssh docker load (%s): %w", name, sshErr)
		}
	}

	return nil
}

func pushToRegistry(cfg config.Config, tags map[string]string) error {
	for name, tag := range tags {
		tagSuffix := tagSuffix(tag)
		target := fmt.Sprintf("%s:%s-%s", cfg.Deploy.Registry.Repository, name, tagSuffix)
		if cfg.Deploy.Registry.RegistryURL != "" {
			target = fmt.Sprintf("%s/%s", strings.TrimSuffix(cfg.Deploy.Registry.RegistryURL, "/"), target)
		}

		tagCmd := exec.Command("docker", "tag", tag, target)
		tagCmd.Stdout = os.Stdout
		tagCmd.Stderr = os.Stderr
		if err := tagCmd.Run(); err != nil {
			return fmt.Errorf("docker tag (%s): %w", name, err)
		}

		pushCmd := exec.Command("docker", "push", target)
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr
		if err := pushCmd.Run(); err != nil {
			return fmt.Errorf("docker push (%s): %w", name, err)
		}
	}

	return nil
}
