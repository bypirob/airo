package docker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"bypirob/airo/src/internal/config"
)

type TagsResult struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func Tags(cfg config.Config, remote bool) ([]string, error) {
	if remote {
		return remoteTags(cfg)
	}
	return localTags(cfg)
}

func localTags(cfg config.Config) ([]string, error) {
	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker images: %w", err)
	}

	repoPrefix := cfg.Name + ":"
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	tags := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" || !strings.HasPrefix(line, repoPrefix) {
			continue
		}
		tag := strings.TrimPrefix(line, repoPrefix)
		if tag == "" || tag == "<none>" {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func remoteTags(cfg config.Config) ([]string, error) {
	if cfg.Deploy.Registry.Repository == "" {
		return nil, fmt.Errorf("deploy.registry.repository is required for remote tags")
	}
	if cfg.Deploy.Registry.RegistryURL == "" {
		return nil, fmt.Errorf("deploy.registry.registry_url is required for remote tags")
	}

	base := strings.TrimSuffix(cfg.Deploy.Registry.RegistryURL, "/")
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "https://" + base
	}

	url := fmt.Sprintf("%s/v2/%s/tags/list", base, cfg.Deploy.Registry.Repository)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch tags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("registry tags request failed: %s", resp.Status)
	}

	var result TagsResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("parse tags response: %w", err)
	}

	return result.Tags, nil
}
