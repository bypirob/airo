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

var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	configDir := flag.String("config", ".", "Path to config directory (default: current directory)")

	commands := []string{"init", "deploy", "compose", "caddy", "version"}
	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand: ", strings.Join(commands, ", "))
		os.Exit(1)
	}

	if !slices.Contains(commands, os.Args[1]) {
		fmt.Println("Expected subcommand: ", strings.Join(commands, ", "))
		os.Exit(1)
	}

	// Handle version command early (doesn't need config or SSH)
	if os.Args[1] == "version" {
		printVersion()
		return
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

	command := os.Args[1]
	user := config.User
	if command == "init" {
		user = "root"
	}
	client, err := goph.New(user, server, auth)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	fmt.Println("✅ SSH client initialized successfully")

	switch command {
	case "init":
		initializeServer(client, config)
	case "deploy":
		serviceToDeploy := "all"
		remainingArgs := flag.Args()
		if len(remainingArgs) > 0 {
			serviceToDeploy = remainingArgs[0]
		}

		serviceNames := []string{}
		servicesToDeploy := config.Services

		if serviceToDeploy != "all" {
			found := false
			for _, service := range config.Services {
				if service.Name == serviceToDeploy {
					servicesToDeploy = []lib.Service{service}
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Error: service '%s' not found in config\n", serviceToDeploy)
				fmt.Println("Available services:")
				for _, service := range config.Services {
					fmt.Printf("  - %s\n", service.Name)
				}
				os.Exit(1)
			}
			fmt.Printf("🚀 Deploying service: %s\n", serviceToDeploy)
		} else {
			fmt.Printf("🚀 Deploying all services\n")
		}

		for _, service := range servicesToDeploy {
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

func printVersion() {
	fmt.Printf("Airo version %s\n", version)
	fmt.Printf("  commit: %s\n", commit)
	fmt.Printf("  built: %s\n", buildDate)
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
	fmt.Println("🔄 Copying compose file")
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
	fmt.Println("✅ Compose file copied and ran successfully")
}

func copyCaddyfileAndReload(client *goph.Client, configDir string) {
	caddyFile := filepath.Join(configDir, "Caddyfile")
	err := client.Upload(caddyFile, "/opt/caddy/Caddyfile")
	if err != nil {
		fmt.Printf("Error copying Caddyfile: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Caddyfile copied successfully")
	_, err = client.Run("docker exec caddy caddy reload --config /etc/caddy/Caddyfile")
	if err != nil {
		fmt.Printf("Error reloading Caddyfile: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Caddyfile reloaded successfully")
}

func initializeServer(client *goph.Client, config lib.Config) error {
	fmt.Println("=== Starting server initialization ===")

	fmt.Println("\nUpdating package lists...")
	out, err := client.Run("sudo apt-get update")
	if err != nil {
		return fmt.Errorf("error updating package lists: %w\n%s", err, string(out))
	}
	fmt.Println("✓ Package lists updated")

	fmt.Println("\nUpgrading existing packages (this may take a while)...")
	out, err = client.Run("sudo DEBIAN_FRONTEND=noninteractive apt-get upgrade -y")
	if err != nil {
		return fmt.Errorf("error upgrading packages: %w\n%s", err, string(out))
	}
	fmt.Println("✓ Packages upgraded")

	fmt.Println("\nInstalling prerequisites...")
	out, err = client.Run("sudo apt-get install -y ca-certificates curl gnupg lsb-release")
	if err != nil {
		return fmt.Errorf("error installing prerequisites: %w\n%s", err, string(out))
	}
	fmt.Println("✓ Prerequisites installed")

	fmt.Println("\nInstalling Docker (this may take a while)...")
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

	fmt.Println("\nCreating deployment user...")
	createUserScript := fmt.Sprintf(`
		# Check if user exists
		if id "%s" &>/dev/null; then
			echo "User %s already exists"
		else
			# Create user with home directory
			sudo useradd -m -s /bin/bash %s
			echo "User %s created successfully"
		fi
		mkdir -p /opt/caddy
		chown -R %s:%s /opt/caddy
	`, config.User, config.User, config.User, config.User, config.User, config.User)
	out, err = client.Run(createUserScript)
	if err != nil {
		return fmt.Errorf("error creating deployment user: %w\n%s", err, string(out))
	}
	fmt.Printf("✓ Deployment user '%s' ready\n", config.User)

	fmt.Println("\nAdding users to docker group...")
	out, err = client.Run(fmt.Sprintf("sudo usermod -aG docker $USER && sudo usermod -aG docker %s", config.User))
	if err != nil {
		return fmt.Errorf("error adding users to docker group: %w\n%s", err, string(out))
	}
	fmt.Println("✓ Users added to docker group")

	fmt.Println("\nCopying SSH authorized keys to deployment user...")
	copyKeysScript := fmt.Sprintf(`
		# Create .ssh directory for the new user
		sudo mkdir -p /home/%s/.ssh
		# Copy root's authorized_keys to the new user
		sudo cp /root/.ssh/authorized_keys /home/%s/.ssh/authorized_keys
		# Set proper ownership and permissions
		sudo chown -R %s:%s /home/%s/.ssh
		sudo chmod 700 /home/%s/.ssh
		sudo chmod 600 /home/%s/.ssh/authorized_keys
	`, config.User, config.User, config.User, config.User, config.User, config.User, config.User)
	out, err = client.Run(copyKeysScript)
	if err != nil {
		return fmt.Errorf("error copying SSH keys: %w\n%s", err, string(out))
	}
	fmt.Printf("✓ SSH keys copied to '%s'\n", config.User)

	fmt.Println("\nVerifying Docker installation...")
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
