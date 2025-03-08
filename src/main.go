package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/bypirob/airo/src/lib"
	"github.com/melbahja/goph"
)

func main() {
	commands := []string{"deploy", "compose", "caddy"}
	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand: ", strings.Join(commands, ", "))
		os.Exit(1)
	}

	if !slices.Contains(commands, os.Args[1]) {
		fmt.Println("Expected subcommand: ", strings.Join(commands, ", "))
		os.Exit(1)
	}

	envFile := "./env.yml"

	config := lib.ReadConfig(envFile)
	fmt.Printf("Config: %+v\n", config)

	server := config.Server

	auth, err := goph.Key(config.SshKey, "")
	if err != nil {
		log.Fatal(err)
	}

	client, err := goph.New(config.User, server, auth)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	if err != nil {
		fmt.Printf("Error initializing SSH client: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SSH client initialized successfully")

	switch os.Args[1] {
	case "deploy":
		serviceNames := []string{}
		for _, service := range config.Services {
			buildDockerImage(service.Build, service.Image)
			pushDockerImage(service.Image)
			pullDockerImage(client, service.Image)
			serviceNames = append(serviceNames, service.Name)
		}

		copyComposeAndRun(client, serviceNames, config.User)
	case "compose":
		copyComposeAndRun(client, []string{}, config.User)
	case "caddy":
		copyCaddyfileAndReload(client)
	default:
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}

func buildDockerImage(folder string, tag string) {
	cmd := exec.Command("docker", "build", "--platform", "linux/amd64", "-t", tag, folder)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error building Docker image: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Docker image built successfully")
}

func pushDockerImage(tag string) {
	cmd := exec.Command("docker", "push", tag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error push Docker image: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Docker image pushed successfully")
}

func pullDockerImage(client *goph.Client, tag string) {
	_, err := client.Run("docker pull " + tag)

	if err != nil {
		fmt.Printf("Error pull Docker image: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Docker image pull successfully")
}

func copyComposeAndRun(client *goph.Client, services []string, user string) {
	fmt.Println("Copying compose file")
	err := client.Upload("./compose.yml", "/home/"+user+"/compose.yml")
	if err != nil {
		fmt.Printf("Error copying compose file: %v\n", err)
		os.Exit(1)
	}

	out, err := client.Run("docker compose up -d " + strings.Join(services, " "))
	if err != nil {
		fmt.Printf("Error running compose: %v\n", err)
		fmt.Println(string(out))
		os.Exit(1)
	}
	fmt.Println("Compose file copied and ran successfully")
}

func copyCaddyfileAndReload(client *goph.Client) {
	client.Upload("./Caddyfile", "/opt/caddy/Caddyfile")
	client.Run("docker exec caddy caddy reload --config /etc/caddy/Caddyfile")
}
