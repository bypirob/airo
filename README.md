# 🚀 Airo

**Deploy your projects directly from your local computer to your production server (VPS) easily.**

Airo builds Docker images and deploys them over SSH or via a registry, driven by `airo.yaml`.

## Why Airo?

Deploying side-projects doesn't have to be complicated or expensive. Kubernetes, Platform as a Service (PaaS) and CI/CD pipelines are a powerful and exciting solutions, but sometimes they're more complex than your project requires. If you enjoy managing your server, it can be significantly cheaper and offer greater control over the technical details.

I want to automate this process and deploy easily to my own server. That's why I've created **Airo**:

- 🚀 **Focus on building your product**, not managing infrastructure.
- 🐳 **Build and deliver Docker images via a registry or direct copy**.
- ⚡️ **Deploy instantly** with a single command from your computer.
- 🔑 **Easily update configurations and containers securely** using SSH.

## Installation

### From Source

```bash
git clone https://github.com/bypirob/airo.git
cd airo
make install
```

## Usage

`airo release` builds, pushes, and deploys in one step, and generates a tag suffix automatically when `--tag` is omitted.

### Configure airo.yaml

```yaml
images:
  app:
    base_image: node:24-alpine
    target_arch: linux/amd64
deploy:
  type: ssh # or registry
  containers:
    - name: "app"
      image: "app"
      port: 3000
      app_port: 3000
      env_file: "/etc/airo/app.env"
      networks:
        - "frontend"
        - "backend"
  ssh:
    host: "192.168.1.100"
    user: "admin"
    port: 22
    identity_file: "~/.ssh/id_rsa"
  registry:
    registry_url: "registry.example.com"
    repository: "my-app"
```

### Commands

```bash
airo build --tag dev --context .
airo push dev
airo deploy --tag dev
airo status
airo tags
airo tags --remote
airo release --tag dev --context .
airo version
```

### Project and config paths

By default, airo reads `airo.yaml` from the current directory. You can point to a different project root or config file:

```bash
airo build --project /path/to/project --config airo.yaml
```
