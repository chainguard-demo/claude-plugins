# Chainguard Plugins for Claude Code

Official Chainguard plugins for Claude Code - bringing supply chain security and secure container images to your development workflow.

## Available Plugins

### 🔍 [chainguard-docs](./chainguard-docs)

Access Chainguard's complete documentation library through a Model Context Protocol (MCP) server.

**Features:**
- Search across all Chainguard documentation
- Get detailed container image information
- Access security guides and CVE management docs
- Learn about Wolfi, apko, melange, and chainctl
- Real-time documentation updates

**Perfect for:**
- Learning about Chainguard products and services
- Finding the right container image for your project
- Understanding supply chain security concepts
- Looking up tool documentation

### 🔧 [chainguard-codegen](./chainguard-codegen)

Generate secure code and configurations using Chainguard Images and tools.

**Features:**
- Generate secure Dockerfiles with Chainguard Images
- Migrate existing Dockerfiles to Chainguard
- Create apko/melange configurations
- Apply security best practices automatically
- Multi-stage build optimization

**Perfect for:**
- Creating new containerized applications
- Migrating existing projects to Chainguard
- Building custom Wolfi-based images
- Implementing security best practices

## Installation

### Add the Chainguard Marketplace

```bash
/plugin marketplace add github.com/chainguard-demo/claude-plugins
```

### Install Both Plugins

```bash
/plugin install chainguard-docs@chainguard-plugins
/plugin install chainguard-codegen@chainguard-plugins
```

### Or Install Individually

**Documentation only:**
```bash
/plugin install chainguard-docs@chainguard-plugins
```

**Code generation only:**
```bash
/plugin install chainguard-codegen@chainguard-plugins
```

## Quick Start

Once installed, try these example prompts:

### Using chainguard-docs

```
"What Chainguard images are available for Python?"
"How does Chainguard handle CVE management?"
"Show me the nginx image documentation"
"Search Chainguard docs for SBOM generation"
```

### Using chainguard-codegen

```
"Generate a Dockerfile for a Python Flask app using Chainguard Images"
"Migrate this Dockerfile to use Chainguard Images"
"Create an apko config for a Node.js application"
"Generate a secure multi-stage Dockerfile for Go"
```

## Requirements

### chainguard-docs
- Docker installed and running (for MCP server)
- Internet connection

### chainguard-codegen
- No additional requirements

## Plugin Architecture

### chainguard-docs

This plugin wraps Chainguard's MCP server (`ghcr.io/chainguard-dev/ai-docs:latest`) which provides:
- Full-text search across documentation
- Structured access to container image docs
- Security and tool reference materials
- Weekly automated updates

### chainguard-codegen

This plugin includes specialized AI skills that help Claude Code:
- Generate secure Dockerfiles following best practices
- Migrate legacy Dockerfiles to Chainguard Images
- Create apko configurations for custom Wolfi images
- Apply security defaults and multi-stage builds

## What are Chainguard Images?

Chainguard Images are secure, minimal container images designed for production use:

- **Zero CVEs**: No known vulnerabilities, regularly updated
- **Minimal**: Only includes what's necessary, reducing attack surface
- **SBOM**: Complete software bill of materials included
- **Signed**: Cryptographically signed with Sigstore
- **Non-root**: Run as non-root user by default
- **Variants**: Production (`:latest`) and development (`:latest-dev`) variants

Learn more at [images.chainguard.dev](https://images.chainguard.dev)

## Examples

### Example 1: Learning About Images

```
You: What Chainguard images are available for Python development?

Claude Code: [Uses chainguard-docs MCP server]
Based on Chainguard documentation, here are the Python images:

- cgr.dev/chainguard/python:latest - Minimal production runtime
- cgr.dev/chainguard/python:latest-dev - Development variant with pip and build tools
...
```

### Example 2: Generating a Dockerfile

```
You: Generate a secure Dockerfile for a Django application

Claude Code: [Uses chainguard-codegen skills]
I'll create a secure multi-stage Dockerfile using Chainguard's Python image:

FROM cgr.dev/chainguard/python:latest-dev AS builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt --user

FROM cgr.dev/chainguard/python:latest
WORKDIR /app
COPY --from=builder /home/nonroot/.local /home/nonroot/.local
COPY --chown=nonroot:nonroot . .
ENV PATH=/home/nonroot/.local/bin:$PATH
EXPOSE 8000
CMD ["python", "manage.py", "runserver", "0.0.0.0:8000"]
```

### Example 3: Migrating Existing Code

```
You: Can you migrate this Dockerfile to use Chainguard Images?

FROM python:3.12-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
CMD ["python", "app.py"]

Claude Code: [Uses chainguard-codegen dockerfile-migrator skill]
I'll migrate this to use Chainguard's Python image with a multi-stage build for better security...
[Shows migrated Dockerfile with explanation of changes]
```

## Documentation

- [Chainguard Academy](https://edu.chainguard.dev) - Tutorials and guides
- [Chainguard Images](https://images.chainguard.dev) - Full image catalog
- [Claude Code Plugins](https://code.claude.com/docs/en/plugins) - Plugin documentation

## Support

- **Documentation**: [edu.chainguard.dev](https://edu.chainguard.dev)
- **Support Portal**: [support.chainguard.dev](https://support.chainguard.dev)
- **Community Slack**: [go.chainguard.dev/slack](https://go.chainguard.dev/slack)
- **GitHub Issues**: [github.com/chainguard-dev/edu/issues](https://github.com/chainguard-dev/edu/issues)

## Contributing

These plugins are maintained by Chainguard. For issues, suggestions, or contributions:

1. Open an issue in the [edu repository](https://github.com/chainguard-dev/edu/issues)
2. Join our [Community Slack](https://go.chainguard.dev/slack)
3. Contact [support@chainguard.dev](mailto:support@chainguard.dev)

## License

Apache-2.0

## About Chainguard

Chainguard is building a new generation of secure software supply chain tooling. We provide:

- **Chainguard Images**: Secure, minimal container images
- **Wolfi**: Undistro for building containers
- **apko**: Build OCI images from APK packages
- **melange**: Build APK packages from source
- **Chainguard Enforce**: Software supply chain security platform

Learn more at [chainguard.dev](https://chainguard.dev)
