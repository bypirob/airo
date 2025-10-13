package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bypirob/airo/src/lib"
	"github.com/melbahja/goph"
)

func main() {
	configDir := flag.String("config", ".", "Path to config directory (default: current directory)")

	commands := []string{"init", "deploy", "compose", "caddy"}
	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand: ", strings.Join(commands, ", "))
		os.Exit(1)
	}

	if !slices.Contains(commands, os.Args[1]) {
		fmt.Println("Expected subcommand: ", strings.Join(commands, ", "))
		os.Exit(1)
	}

	flag.CommandLine.Parse(os.Args[2:])

	envFile := filepath.Join(*configDir, "env.yml")

	config := lib.ReadConfig(envFile)
	fmt.Printf("✅ Config loaded successfully\n")

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

	fmt.Println("✅ SSH client initialized successfully")

	switch os.Args[1] {
	case "init":
		initializeServer(client)
	case "deploy":
		serviceNames := []string{}
		for _, service := range config.Services {
			buildDockerImage(service.Build, service.Image)
			if config.Transport == "copy" {
				if err := copyImageViaSSH(client, service.Image, config.SshKey, config.User, config.Server); err != nil {
					fmt.Printf("Error copying image via SSH: %v\n", err)
					os.Exit(1)
				}
			} else {
				pushDockerImage(service.Image)
				pullDockerImage(client, service.Image)
			}
			serviceNames = append(serviceNames, service.Name)
		}

		copyComposeAndRun(client, serviceNames, config.User, *configDir)
	case "compose":
		copyComposeAndRun(client, []string{}, config.User, *configDir)
	case "caddy":
		copyCaddyfileAndReload(client, *configDir)
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

func copyComposeAndRun(client *goph.Client, services []string, user string, configDir string) {
	fmt.Println("Copying compose file")
	composeFile := filepath.Join(configDir, "compose.yml")
	path := "/home/" + user + "/compose.yml"
	if user == "root" {
		path = "/root/compose.yml"
	}
	err := client.Upload(composeFile, path)
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

func copyCaddyfileAndReload(client *goph.Client, configDir string) {
	caddyFile := filepath.Join(configDir, "Caddyfile")
	client.Upload(caddyFile, "/opt/caddy/Caddyfile")
	client.Run("docker exec caddy caddy reload --config /etc/caddy/Caddyfile")
}

func initializeServer(client *goph.Client) error {
	fmt.Println("=== Starting server initialization ===")

	fmt.Println("\n[1/6] Updating package lists...")
	out, err := client.Run("sudo apt-get update")
	if err != nil {
		return fmt.Errorf("error updating package lists: %w\n%s", err, string(out))
	}
	fmt.Println("✓ Package lists updated")

	fmt.Println("\n[2/6] Upgrading existing packages (this may take a while)...")
	out, err = client.Run("sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y")
	if err != nil {
		return fmt.Errorf("error upgrading packages: %w\n%s", err, string(out))
	}
	fmt.Println("✓ Packages upgraded")

	fmt.Println("\n[3/6] Installing prerequisites...")
	out, err = client.Run("sudo apt-get install -y ca-certificates curl gnupg lsb-release")
	if err != nil {
		return fmt.Errorf("error installing prerequisites: %w\n%s", err, string(out))
	}
	fmt.Println("✓ Prerequisites installed")

	fmt.Println("\n[4/6] Installing Docker (this may take a while)...")
	dockerInstallScript := `
		# Add Docker's official GPG key
		sudo install -m 0755 -d /etc/apt/keyrings
		curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor --yes -o /etc/apt/keyrings/docker.gpg
		sudo chmod a+r /etc/apt/keyrings/docker.gpg
		
		# Set up Docker repository
		echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
		
		# Install Docker
		sudo apt-get update
		sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
	`
	out, err = client.Run(dockerInstallScript)
	if err != nil {
		return fmt.Errorf("error installing Docker: %w\n%s", err, string(out))
	}
	fmt.Println("✓ Docker installed")

	fmt.Println("\n[5/6] Adding user to docker group...")
	out, err = client.Run("sudo usermod -aG docker $USER")
	if err != nil {
		return fmt.Errorf("error adding user to docker group: %w\n%s", err, string(out))
	}
	fmt.Println("✓ User added to docker group")

	fmt.Println("\n[6/6] Verifying Docker installation...")
	out, err = client.Run("sudo docker --version && sudo docker compose version")
	if err != nil {
		return fmt.Errorf("error verifying Docker installation: %w\n%s", err, string(out))
	}
	fmt.Printf("✓ Docker verified:\n%s\n", string(out))

	fmt.Println("\n=== Server initialization complete! ===")
	fmt.Println("\nNote: You may need to log out and back in for docker group permissions to take effect.")
	return nil
}

func copyImageViaSSH(client *goph.Client, tag string, sshKey string, user string, server string) error {
	fmt.Printf("Copying image %s via SCP...\n", tag)

	// 1. Save image to disk
	localTar := "/tmp/" + strings.ReplaceAll(tag, "/", "_") + ".tar"
	fmt.Printf("Saving image to %s...\n", localTar)

	saveCmd := exec.Command("docker", "save", "-o", localTar, tag)
	saveCmd.Stdout = os.Stdout
	saveCmd.Stderr = os.Stderr
	if err := saveCmd.Run(); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	defer os.Remove(localTar) // Clean up local tar file

	fmt.Println("Image saved successfully")

	// 2. Copy to server using SCP
	remoteTar := "/tmp/" + filepath.Base(localTar)

	fmt.Printf("Copying to server via SCP...\n")
	scpCmd := exec.Command("scp", "-i", sshKey, localTar, user+"@"+server+":"+remoteTar)
	scpCmd.Stdout = os.Stdout
	scpCmd.Stderr = os.Stderr
	if err := scpCmd.Run(); err != nil {
		return fmt.Errorf("failed to copy image via SCP: %w", err)
	}

	fmt.Println("Image copied to server successfully")

	// 3. Load image on server
	fmt.Println("Loading image on server...")
	_, err := client.Run("docker load -i " + remoteTar + " && rm " + remoteTar)
	if err != nil {
		return fmt.Errorf("failed to load image on server: %w", err)
	}

	fmt.Printf("Image %s loaded successfully on server\n", tag)
	return nil
}
