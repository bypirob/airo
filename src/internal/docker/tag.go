package docker

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"bypirob/airo/src/internal/config"
)

func resolveTags(cfg config.Config, projectPath, tag string) (map[string]string, error) {
	if len(cfg.Images) == 0 {
		return nil, fmt.Errorf("images is required")
	}
	if projectPath == "" {
		projectPath = "."
	}

	if tag == "" {
		suffix, err := defaultTagSuffix(projectPath)
		if err != nil {
			return nil, err
		}
		return imageTags(cfg, suffix), nil
	}

	if strings.Contains(tag, ":") {
		if len(cfg.Images) != 1 {
			return nil, fmt.Errorf("tag with repository is only supported with a single image")
		}
		tags := make(map[string]string, len(cfg.Images))
		for name := range cfg.Images {
			tags[name] = tag
		}
		return tags, nil
	}

	return imageTags(cfg, tag), nil
}

func imageTags(cfg config.Config, suffix string) map[string]string {
	tags := make(map[string]string, len(cfg.Images))
	for name := range cfg.Images {
		tags[name] = imageTag(cfg, name, suffix)
	}
	return tags
}

func imageTag(cfg config.Config, imageName, suffix string) string {
	return fmt.Sprintf("%s:%s", imageName, suffix)
}

func tagSuffix(tag string) string {
	if idx := strings.LastIndex(tag, ":"); idx != -1 {
		return tag[idx+1:]
	}
	return tag
}

func defaultTagSuffix(projectPath string) (string, error) {
	sha, err := gitShortSHA(projectPath)
	if err != nil {
		return "", err
	}

	timestamp := time.Now().UTC().Format("20060102-1504")
	return fmt.Sprintf("%s-%s", timestamp, sha), nil
}

func DefaultTagSuffix(projectPath string) (string, error) {
	return defaultTagSuffix(projectPath)
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
