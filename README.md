# 🚀 Airo

**Deploy your projects directly from your local computer to your production server easily.**

Airo helps you deploying containers to your self-hosted server, without worrying about configuring pipelines, serverless services or different platforms. Just your self-hosted servers.

## Why Airo?

Deploying side-projects doesn't have to be complicated or expensive. Kubernetes, Platform as a Service (PaaS) and CI/CD pipelines are a powerful and exciting solutions, but sometimes they're more complex than your project requires. If you enjoy managing your server, it can be significantly cheaper and offer greater control over the technical details.

I want to automate this process and deploy easily to my own server. That's why I've created **Airo**:

- 🚀 **Focus on building your product**, not managing infrastructure.
- 🐳 **Build and push Docker images directly from your local machine to a container registry**.
- ⚡️ **Deploy instantly** with a single command from your computer.
- 🔑 **Easily update configurations and containers securely** using SSH.
- 🌐 **Set up HTTPS and reverse proxy automatically** using Caddy.

## How It Works

Deploying with Airo is easy:

1. Define your services in a `compose.yml` file.
2. Configure your deployment details in `env.yml` (server details, Docker images, etc.).
3. Prepare your Dockerfile.
4. Set up your Caddyfile for automatic HTTPS and reverse proxy.

After this initial setup, deploying new updates is just a simple command away:

```bash
airo deploy
```

## Installation

### From Source

```bash
git clone https://github.com/yourusername/airo.git
cd airo
make install
airo deploy
```

## Usage

1. Create a new project directory and navigate to it:
   ```bash
   mkdir my-project/.deploy
   cd my-project/.deploy
   ```

2. Configure your `env.yaml` file:
   ```yaml
   server: your-server-ip
   user: your-ssh-user
   ssh_key: /path/to/your/ssh/key
   # add all your services that will be built here
   services:
     - name: nextjs
       image: your-registry/api:latest
       build: /home/user/your-project
   ```

3. Add a Dockerfile to your project:
    ```dockerfile
    FROM node:20-alpine
    WORKDIR /app
    COPY . .
    RUN npm install
    ```

4. Configure your `compose.yml` file with your services:
   ```yaml
   services:
      postgres:
        image: postgres:latest
        restart: unless-stopped
        ports:
          - '5432:5432'
        volumes:
          - ./postgres_data:/var/lib/postgresql/data
        environment:
          POSTGRES_PASSWORD: some-temporary-password
     nextjs:
       image: your-registry/front:latest
       restart: unless-stopped
       ports:
         - '3000:3000'
       volumes:
         - ./.env:/app/.env
   ```

5. Deploy your project:
   ```bash
   airo deploy
   ```

6. If you want to deploy your project without building the docker image, you can use the `compose` subcommand:
  ```bash
  airo compose
  ```

7. Also if you want to update your Caddyfile, you can use the `caddy` subcommand:
  ```bash
  airo caddy
  ```

