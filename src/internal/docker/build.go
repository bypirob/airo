package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"bypirob/airo/src/internal/config"
)

func BuildImage(cfg config.Config, projectPath, tag, contextPath string) error {
	if contextPath == "" {
		contextPath = "."
	}
	if projectPath == "" {
		projectPath = "."
	}

	tags, err := resolveTags(cfg, projectPath, tag)
	if err != nil {
		return err
	}

	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	if !filepath.IsAbs(contextPath) {
		contextPath = filepath.Join(projectPath, contextPath)
	}

	for name, image := range cfg.Images {
		imageTag := tags[name]
		args := []string{
			"buildx", "build",
			"--platform", image.TargetArch,
			"--tag", imageTag,
			"--file", dockerfilePath,
			contextPath,
		}

		cmd := exec.Command("docker", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("docker buildx build (%s): %w", name, err)
		}
	}

	return nil
}
