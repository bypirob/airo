# 🚀 Airo

**Effortless deployment for your side-projects**.

Airo simplifies deploying your side-projects from your local computer to a self-hosted server using Docker, Docker Compose, SSH and Caddyfile. 

Airo helps you focus on building your product, not managing infrastructure. Kubernetes or a CI/CD pipeline are cool,
but sometimes I just want to keep things simple in a single server.

## Key features

- 🐳 Build and push your project docker images, directly from your local environment to a docker registry.
- ⚡️ Airo connects to your server via SSH for updating your configs and deploying your containers.
- 🪄 All you need is docker, docker-compose, ssh and a Caddyfile.

## How it works

The idea is to have a simple way to deploy your project to a self-hosted server.

First, the initial setup is done manually:

1. Write a compose.yml file that defines all your docker services.
2. Configure your env.yml file set your server, and the images that will be built.
3. Add a Dockerfile to build your project.
4. Add a Caddyfile that will be used as a reverse proxy for your services.

Then, you can use the `airo deploy` command to deploy your project to your server every
time you finish working on a feature.

> You need to have a server running somewhere. It involves configuring some options manually right
> now that I would like to automate.

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

4. Deploy your project:
   ```bash
   airo deploy
   ```

If you want to deploy your project without building the docker image, you can use the `compose` subcommand:
```bash
airo compose
```

Also if you want to update your Caddyfile, you can use the `caddy` subcommand:
```bash
airo caddy
```
