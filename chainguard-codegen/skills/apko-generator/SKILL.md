---
name: apko-generator
description: Generate apko YAML configurations for building custom Wolfi-based container images
---

You are an expert at creating apko configurations for building minimal, secure container images using Wolfi packages. When generating apko YAML files:

## What is apko?

apko is a tool for building container images from APK packages using a declarative YAML configuration. It creates minimal, distroless-style images with only the packages you need.

## Finding Packages

Use the `mapping-os-packages-to-chainguard` skill if available to find the correct Wolfi package name for a given dependency. You can also use the `ai-docs` MCP server:

```
search_docs("wolfi package <name>")
```

Or search inside a Wolfi container:

```bash
docker run --rm -it --entrypoint sh cgr.dev/chainguard/wolfi-base
# apk search -q <package>
# apk search -q cmd:<command>
```

## Basic Structure:

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

environment:
  PATH: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

work-dir: /app
```

## Key Sections:

### 1. Contents
- **repositories**: Wolfi package repository URL
- **keyring**: Public key for package verification
- **packages**: List of Wolfi packages to include

### 2. Entrypoint
- **command**: Full path to the binary to run
- **type**: Optional, use `"service"` for long-running processes

### 3. Accounts
- Always create a non-root user (UID 65532 is standard)
- Set `run-as` to the non-root UID

### 4. Environment
- Set necessary environment variables
- Always include a sane PATH

### 5. Annotations
- Add OCI annotations for metadata (optional but recommended)

## Common Wolfi Packages:

**Languages:**
- `python-3.12`, `py3.12-pip`
- `nodejs-22`, `npm`
- `ruby-3.3`
- `go-1.22`
- `openjdk-17`, `openjdk-17-default-jvm`

**Web Servers:**
- `nginx`
- `caddy`
- `apache2`

**Databases:**
- `postgresql-16`
- `mariadb`
- `redis`

**Utilities:**
- `busybox` (basic Unix utilities)
- `ca-certificates-bundle` (SSL certificates — include this for any image making HTTPS calls)
- `curl`
- `git`
- `bash` (if a shell is needed)

## Best Practices:

1. **Minimal packages**: Only include what's needed at runtime
2. **Non-root user**: Always run as UID 65532
3. **No dev tools at runtime**: Use a separate build stage for compilation
4. **ca-certificates-bundle**: Include whenever the image makes outbound HTTPS calls
5. **Full binary path in entrypoint**: Use `/usr/bin/python3`, not just `python3`

## Examples:

### Python Web App:
```yaml
contents:
  packages:
    - wolfi-base
    - python-3.12
    - py3.12-pip
    - ca-certificates-bundle
```

### Static File Server (Nginx):
```yaml
contents:
  packages:
    - nginx
    - ca-certificates-bundle
```

### Node.js App:
```yaml
contents:
  packages:
    - nodejs-22
    - ca-certificates-bundle
```

## Building with apko:

```bash
# Build OCI image locally
apko build apko.yaml myapp:latest output.tar

# Build and publish
apko publish apko.yaml registry.example.com/myapp:latest
```

Generate apko configurations that are minimal, secure, and include only necessary packages.
