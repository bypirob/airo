package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bypirob/airo/src/internal/config"
)

func PushImage(cfg config.Config, projectPath, tag string) error {
	if tag == "" {
		defaultTag, err := DefaultTag(cfg, projectPath)
		if err != nil {
			return err
		}
		tag = defaultTag
	}
	if projectPath == "" {
		projectPath = "."
	}

	switch cfg.Deploy.Type {
	case "ssh":
		return pushOverSSH(cfg, tag)
	case "registry":
		return pushToRegistry(cfg, tag)
	default:
		return fmt.Errorf("unsupported deploy.type %q", cfg.Deploy.Type)
	}
}

func pushOverSSH(cfg config.Config, tag string) error {
	saveCmd := exec.Command("docker", "save", tag)
	sshCmd := sshCommand(cfg, "docker", "load")

	pipe, err := saveCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("prepare docker save output: %w", err)
	}

	saveCmd.Stderr = os.Stderr
	sshCmd.Stdin = pipe
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Start(); err != nil {
		return fmt.Errorf("start ssh: %w", err)
	}
	if err := saveCmd.Start(); err != nil {
		return fmt.Errorf("start docker save: %w", err)
	}

	saveErr := saveCmd.Wait()
	sshErr := sshCmd.Wait()

	if saveErr != nil {
		return fmt.Errorf("docker save: %w", saveErr)
	}
	if sshErr != nil {
		return fmt.Errorf("ssh docker load: %w", sshErr)
	}

	return nil
}

func pushToRegistry(cfg config.Config, tag string) error {
	tagSuffix := tag
	if idx := strings.LastIndex(tag, ":"); idx != -1 {
		tagSuffix = tag[idx+1:]
	}

	target := fmt.Sprintf("%s:%s", cfg.Deploy.Registry.Repository, tagSuffix)
	if cfg.Deploy.Registry.RegistryURL != "" {
		target = fmt.Sprintf("%s/%s", strings.TrimSuffix(cfg.Deploy.Registry.RegistryURL, "/"), target)
	}

	tagCmd := exec.Command("docker", "tag", tag, target)
	tagCmd.Stdout = os.Stdout
	tagCmd.Stderr = os.Stderr
	if err := tagCmd.Run(); err != nil {
		return fmt.Errorf("docker tag: %w", err)
	}

	pushCmd := exec.Command("docker", "push", target)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("docker push: %w", err)
	}

	return nil
}
