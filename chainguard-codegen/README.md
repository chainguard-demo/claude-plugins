# Chainguard Code Generation Plugin

Generate secure Dockerfiles, apko/melange configurations, and migrate existing projects to use Chainguard Images.

## Features

- **Dockerfile Generation**: Create secure, minimal Dockerfiles using Chainguard Images
- **Dockerfile Migration**: Convert existing Dockerfiles to use Chainguard Images with the [dfc](https://github.com/chainguard-dev/dfc) (Dockerfile Converter) MCP server
- **apko Configurations**: Generate apko YAML files for custom Wolfi-based images
- **Best Practices**: Automatic application of security and container best practices
- **Multi-stage Builds**: Proper separation of build and runtime environments
- **Intelligent Mapping**: Uses Chainguard's built-in image and package mappings

## Integrated Tools

This plugin includes:
- **dfc MCP Server**: Automated Dockerfile conversion using Chainguard's mapping database
- **AI Skills**: Expert guidance for generation, migration, and best practices
- **Built-in Mappings**: Comprehensive image and package conversion tables

## Requirements

- **Docker**: Required to run the dfc MCP server container
- Internet access to pull `cgr.dev/chainguard/dfc-mcp`

## Installation

Users can install this plugin by adding the Chainguard marketplace:

```bash
/plugin marketplace add github.com/chainguard-demo/claude-plugins
/plugin install chainguard-codegen@chainguard-plugins
```

The dfc MCP server container will be automatically pulled from GitHub Container Registry and run when needed. The container is:
- Built on Chainguard's minimal base images
- Signed with Sigstore cosign
- Includes an SBOM (Software Bill of Materials)
- Automatically updated when dfc changes

## Capabilities

### 1. Dockerfile Generation

Ask Claude Code to generate Dockerfiles with Chainguard Images:

```
"Generate a Dockerfile for a Python Flask app using Chainguard Images"
"Create a multi-stage Dockerfile for a Node.js application"
"Generate a Dockerfile for a Go service that produces a static binary"
```

The plugin will:
- Select the appropriate Chainguard base image
- Use multi-stage builds when beneficial
- Set up non-root users correctly
- Apply security best practices
- Minimize image layers and size

### 2. Dockerfile Migration

Migrate existing Dockerfiles to Chainguard Images:

```
"Migrate this Dockerfile to use Chainguard Images"
"Convert my python:3.12-slim Dockerfile to Chainguard"
"Help me switch from Alpine to Chainguard Images"
```

The plugin will:
- Use the dfc MCP server for automated conversion with built-in mappings
- Identify the current base image
- Map to the appropriate Chainguard equivalent
- Convert package manager commands (apt-get → apk)
- Handle permission and path changes
- Adjust for non-root user requirements
- Explain the security benefits

#### dfc MCP Tools Available

When you ask Claude Code to migrate Dockerfiles, it can use these MCP tools:

- **`convert_dockerfile`**: Converts a complete Dockerfile using Chainguard's mapping database
  - Handles FROM line mapping
  - Converts RUN commands with package managers
  - Applies tag mapping logic
  - Maintains comments and formatting

- **`analyze_dockerfile`**: Analyzes Dockerfile structure before migration
  - Identifies package managers used
  - Detects base images
  - Reports multi-stage build structure

### 3. apko Configuration Generation

Create apko YAML files for custom Wolfi-based images:

```
"Generate an apko config for a Python web application"
"Create an apko.yaml for a custom Nginx image"
"Build an apko configuration with PostgreSQL and utilities"
```

The plugin will:
- Select necessary Wolfi packages
- Configure non-root user (UID 65532)
- Set proper entrypoint and environment
- Keep the image minimal and secure

## Chainguard Images Reference

### Common Images

| Use Case | Image | Variants |
|----------|-------|----------|
| Python | `cgr.dev/chainguard/python` | `:latest`, `:latest-dev`, `-fips` |
| Node.js | `cgr.dev/chainguard/node` | `:latest`, `:latest-dev` |
| Go | `cgr.dev/chainguard/go` | `:latest` (builder) |
| Static binaries | `cgr.dev/chainguard/static` | `:latest` (runtime) |
| Nginx | `cgr.dev/chainguard/nginx` | `:latest` |
| PostgreSQL | `cgr.dev/chainguard/postgres` | `:latest` |
| Redis | `cgr.dev/chainguard/redis` | `:latest` |
| Wolfi base | `cgr.dev/chainguard/wolfi-base` | `:latest` |

### Image Variants

- **`:latest`** - Production runtime image (minimal, no shell or package manager)
- **`:latest-dev`** - Development image (includes shell, apk, build tools)
- **`-fips`** - FIPS 140-2 compliant variants

## Best Practices Applied

When generating or migrating code, this plugin ensures:

1. **Multi-stage builds**: Build with `:latest-dev`, run with `:latest`
2. **Non-root user**: All containers run as UID 65532 (nonroot)
3. **Minimal layers**: Commands are combined efficiently
4. **No secrets**: Never includes secrets in images
5. **Proper permissions**: Files are owned by nonroot user
6. **Security defaults**: Following Chainguard and Docker best practices

## Example Outputs

### Simple Python Dockerfile

```dockerfile
FROM cgr.dev/chainguard/python:latest-dev AS builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt --user

FROM cgr.dev/chainguard/python:latest
WORKDIR /app
COPY --from=builder /home/nonroot/.local /home/nonroot/.local
COPY --chown=nonroot:nonroot . .
ENV PATH=/home/nonroot/.local/bin:$PATH
CMD ["python", "app.py"]
```

### apko Configuration

```yaml
contents:
  repositories:
    - https://packages.wolfi.dev/os
  keyring:
    - https://packages.wolfi.dev/os/wolfi-signing.rsa.pub
  packages:
    - wolfi-base
    - python-3.12
    - py3.12-pip

entrypoint:
  command: /usr/bin/python3

accounts:
  groups:
    - groupname: nonroot
      gid: 65532
  users:
    - username: nonroot
      uid: 65532
      gid: 65532
  run-as: 65532

work-dir: /app
```

## Skills Included

### Plugin Skills (automatically active when plugin is installed)

- **dockerfile-generator**: Expert at creating secure Dockerfiles
- **dockerfile-migrator**: Migrates existing Dockerfiles to Chainguard
- **apko-generator**: Creates apko configurations for Wolfi images

### Standalone Skills (install separately for deeper migration guidance)

The `skills/` directory includes three additional skills in standard SKILL.md
format, designed for use with the `migrating-dockerfiles-to-chainguard` workflow:

- **migrating-dockerfiles-to-chainguard**: Step-by-step guidance for the full
  migration process — clarifying requirements, converting FROM statements,
  handling packages, users, entrypoints, and troubleshooting
- **mapping-container-images-to-chainguard**: Maps upstream container image
  references to their Chainguard equivalents, including tag rules, FIPS
  variants, and digest handling
- **mapping-os-packages-to-chainguard**: Maps OS package names across Alpine,
  Debian, and Fedora ecosystems to their Chainguard/Wolfi equivalents

To install the standalone skills:

```bash
cp -r chainguard-codegen/skills/migrating-dockerfiles-to-chainguard ~/.claude/skills/
cp -r chainguard-codegen/skills/mapping-container-images-to-chainguard ~/.claude/skills/
cp -r chainguard-codegen/skills/mapping-os-packages-to-chainguard ~/.claude/skills/
```

## Resources

- [Chainguard Images](https://images.chainguard.dev) - Full image catalog
- [apko Documentation](https://edu.chainguard.dev/open-source/apko/)
- [Wolfi Documentation](https://edu.chainguard.dev/open-source/wolfi/)
- [Dockerfile Best Practices](https://edu.chainguard.dev/chainguard/chainguard-images/getting-started/)

## Support

- [Chainguard Academy](https://edu.chainguard.dev)
- [Chainguard Support](https://support.chainguard.dev)
- [Community Slack](https://go.chainguard.dev/slack)
