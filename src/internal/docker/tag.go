package docker

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"bypirob/airo/src/internal/config"
)

func DefaultTag(cfg config.Config, projectPath string) (string, error) {
	if projectPath == "" {
		projectPath = "."
	}

	sha, err := gitShortSHA(projectPath)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format("20060102-1504")
	return fmt.Sprintf("%s:%s-%s", cfg.Name, timestamp, sha), nil
}

func gitShortSHA(projectPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("resolve git commit: %w (%s)", err, strings.TrimSpace(string(output)))
	}

	sha := strings.TrimSpace(string(output))
	if sha == "" {
		return "", fmt.Errorf("resolve git commit: empty sha")
	}

	return sha, nil
}
