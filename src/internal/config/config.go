package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const (
	DefaultName       = "airo"
	DefaultBaseImage  = "node:24-alpine"
	DefaultTargetArch = "linux/amd64"
	DefaultInstallCmd = "npm ci"
	DefaultBuildCmd   = "npm run build"
	DefaultStartCmd   = "npm start"
)

type Config struct {
	Name      string          `yaml:"name"`
	Container ContainerConfig `yaml:"container"`
	Deploy    DeployConfig    `yaml:"deploy"`
}

type ContainerConfig struct {
	BaseImage  string `yaml:"base_image"`
	TargetArch string `yaml:"target_arch"`
	InstallCmd string `yaml:"install_cmd"`
	BuildCmd   string `yaml:"build_cmd"`
	StartCmd   string `yaml:"start_cmd"`
	Port       int    `yaml:"port"`
	AppPort    int    `yaml:"app_port"`
}

type DeployConfig struct {
	Type     string         `yaml:"type"`
	EnvFile  string         `yaml:"env_file"`
	Networks []string       `yaml:"networks"`
	SSH      SSHConfig      `yaml:"ssh"`
	Registry RegistryConfig `yaml:"registry"`

	// Legacy ports under deploy; mapped to container when set.
	Port    int `yaml:"port"`
	AppPort int `yaml:"app_port"`
}

type SSHConfig struct {
	Host         string `yaml:"host"`
	User         string `yaml:"user"`
	Port         int    `yaml:"port"`
	IdentityFile string `yaml:"identity_file"`
}

type RegistryConfig struct {
	RegistryURL string `yaml:"registry_url"`
	Repository  string `yaml:"repository"`
}

func Load(projectPath, configPath string) (Config, error) {
	if projectPath == "" {
		projectPath = "."
	}
	if configPath == "" {
		configPath = "airo.yaml"
	}

	fullPath := filepath.Join(projectPath, configPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", fullPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", fullPath, err)
	}

	applyDefaults(&cfg)
	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Name == "" {
		cfg.Name = DefaultName
	}
	if cfg.Container.BaseImage == "" {
		cfg.Container.BaseImage = DefaultBaseImage
	}
	if cfg.Container.TargetArch == "" {
		cfg.Container.TargetArch = DefaultTargetArch
	}
	if cfg.Container.InstallCmd == "" {
		cfg.Container.InstallCmd = DefaultInstallCmd
	}
	if cfg.Container.BuildCmd == "" {
		cfg.Container.BuildCmd = DefaultBuildCmd
	}
	if cfg.Container.StartCmd == "" {
		cfg.Container.StartCmd = DefaultStartCmd
	}

	if cfg.Container.Port == 0 && cfg.Deploy.Port != 0 {
		cfg.Container.Port = cfg.Deploy.Port
	}
	if cfg.Container.AppPort == 0 && cfg.Deploy.AppPort != 0 {
		cfg.Container.AppPort = cfg.Deploy.AppPort
	}
}

func validate(cfg Config) error {
	if cfg.Deploy.Type == "" {
		return fmt.Errorf("deploy.type is required")
	}
	switch cfg.Deploy.Type {
	case "ssh", "registry":
	default:
		return fmt.Errorf("deploy.type must be ssh or registry")
	}

	if cfg.Deploy.Type == "ssh" && cfg.Deploy.SSH.Host == "" {
		return fmt.Errorf("deploy.ssh.host is required for ssh deploys")
	}
	if cfg.Deploy.Type == "registry" && cfg.Deploy.Registry.Repository == "" {
		return fmt.Errorf("deploy.registry.repository is required for registry deploys")
	}

	if cfg.Container.Port != 0 && cfg.Container.AppPort == 0 {
		return fmt.Errorf("container.app_port is required when container.port is set")
	}

	return nil
}
