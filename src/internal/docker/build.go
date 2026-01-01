package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"bypirob/airo/src/internal/config"
)

func BuildImage(cfg config.Config, projectPath, tag, contextPath string) error {
	if tag == "" {
		defaultTag, err := DefaultTag(cfg, projectPath)
		if err != nil {
			return err
		}
		tag = defaultTag
	}
	if contextPath == "" {
		contextPath = "."
	}
	if projectPath == "" {
		projectPath = "."
	}

	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	if !filepath.IsAbs(contextPath) {
		contextPath = filepath.Join(projectPath, contextPath)
	}

	args := []string{
		"buildx", "build",
		"--platform", cfg.Container.TargetArch,
		"--tag", tag,
		"--file", dockerfilePath,
		contextPath,
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker buildx build: %w", err)
	}

	return nil
}
