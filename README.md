# 🚀 Airo

**Deploy your projects directly from your local computer to your production server easily.**

Airo helps you deploying containers to your self-hosted server, without worrying about configuring pipelines, serverless or services. Just your servers.

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

That's it—simple, repeatable, and hassle-free.

---

**Stop dealing with deployment headaches. Deploy easily with Airo.**

