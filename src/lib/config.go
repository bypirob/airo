package lib

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Service struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
	Build string `yaml:"build"`
}

type Config struct {
	Services  []Service `yaml:"services"`
	Server    string    `yaml:"server"`
	User      string    `yaml:"user"`
	SshKey    string    `yaml:"ssh_key"`
	Transport string    `yaml:"transport"`
}

func ReadConfig(file string) Config {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	if config.Transport == "" {
		config.Transport = "registry"
	}

	if config.User == "" {
		config.User = "airo"
	}

	return config
}
