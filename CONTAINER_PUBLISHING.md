# Container Publishing Setup

This document explains how both Chainguard Claude Code plugins publish their MCP servers as containers.

## Overview

Both plugins use **containerized MCP servers** published to GitHub Container Registry (ghcr.io):

1. **chainguard-docs**: `ghcr.io/chainguard-dev/ai-docs:latest`
2. **chainguard-codegen**: `cgr.dev/chainguard/dfc-mcp`

This provides a consistent, dependency-free experience for users.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                 Claude Code Plugins                         │
│                                                             │
│  chainguard-docs          chainguard-codegen               │
│  └── .mcp.json            └── .mcp.json                    │
│      ↓                        ↓                             │
└──────┼─────────────────────────┼──────────────────────────┘
       │                         │
       ↓                         ↓
  ┌────────────────┐      ┌──────────────────┐
  │ ghcr.io/       │      │ cgr.dev/         │
  │ chainguard-dev/│      │ chainguard/      │
  │ ai-docs:latest │      │ dfc-mcp          │
  └────────────────┘      └──────────────────┘
```

## dfc MCP Server Publishing

### What Was Created

**File**: `dfc/.github/workflows/publish-mcp-server.yaml`

This workflow:
1. Triggers on:
   - Pushes to `main` branch (publishes `:latest`)
   - Release events (publishes versioned tags like `:v1.0.0`)
   - Manual dispatch
2. Builds the container from `mcp-server/Dockerfile`
3. Pushes to `cgr.dev/chainguard/dfc-mcp`
4. Signs with cosign
5. Attaches an SBOM
6. Supports multi-arch (amd64, arm64)

### Container Details

**Base Images**:
- Builder: `cgr.dev/chainguard/go:latest`
- Runtime: `cgr.dev/chainguard/static:latest`

**Features**:
- Minimal distroless runtime
- No shell or package manager
- Signed with Sigstore
- SBOM attached
- Multi-architecture support

### Publishing Triggers

The workflow runs when:

1. **Code changes**: Any change to `mcp-server/**` pushed to main
2. **Workflow changes**: Modifications to the workflow file itself
3. **Releases**: When a new version is released
4. **Manual**: Via GitHub Actions UI

### Tags Published

- `latest` - Latest from main branch
- `v1.2.3` - Semantic version tags
- `v1.2` - Major.minor tags
- `v1` - Major version tags
- `sha-<commit>` - Commit-specific tags for traceability

## Plugin Configuration

### chainguard-docs

```json
{
  "mcpServers": {
    "chainguard-docs": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "ghcr.io/chainguard-dev/ai-docs:latest",
        "serve-mcp"
      ]
    }
  }
}
```

### chainguard-codegen

```json
{
  "mcpServers": {
    "dfc": {
      "command": "docker",
      "args": [
        "run",
        "--rm",
        "-i",
        "cgr.dev/chainguard/dfc-mcp"
      ]
    }
  }
}
```

## User Experience

When users install the plugins:

1. Claude Code reads the `.mcp.json` configuration
2. On first use, the container image is pulled automatically
3. The container runs in stdio mode, communicating via stdin/stdout
4. No additional dependencies needed (Go, Python, etc.)

## Security Features

Both containers include:

- **Minimal base images**: Chainguard's distroless images
- **Signed**: Cryptographic signatures with Sigstore cosign
- **SBOM**: Complete software bill of materials
- **No shell**: Reduces attack surface
- **Non-root**: Run as non-root user by default
- **Verified**: Can be verified with cosign

### Verification Example

```bash
# Verify signature
cosign verify cgr.dev/chainguard/dfc-mcp \
  --certificate-identity "https://github.com/chainguard-images/images/.github/workflows/release.yaml@refs/heads/main" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com

# View SBOM
cosign download sbom cgr.dev/chainguard/dfc-mcp
```

## Development Workflow

### Testing Locally

Before publishing, test the container locally:

```bash
# Build locally
cd dfc/mcp-server
docker build -t dfc-mcp-test .

# Test it works
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | \
  docker run --rm -i dfc-mcp-test
```

### Publishing Process

1. **Make changes** to `mcp-server/` code in dfc repo
2. **Commit and push** to main branch
3. **Workflow runs automatically**
4. **Container published** to ghcr.io
5. **Users get updates** on next plugin use

### Release Process

For versioned releases:

1. **Tag a release** in the dfc repo: `git tag v1.0.0`
2. **Push the tag**: `git push origin v1.0.0`
3. **Workflow publishes** versioned containers
4. **Users can pin versions** if needed

## Troubleshooting

### Container Won't Pull

```bash
# Check if you can pull manually
docker pull cgr.dev/chainguard/dfc-mcp

# Check GitHub Container Registry status
# https://www.githubstatus.com/
```

### Permission Issues

The GitHub Actions workflow needs:
- `packages: write` permission (to push to ghcr.io)
- `id-token: write` permission (for cosign signing)
- `contents: read` permission (to checkout code)

These are already configured in the workflow file.

### Build Failures

Check:
1. GitHub Actions logs in the dfc repository
2. Dockerfile syntax is valid
3. Go dependencies are correct in `go.mod`
4. Base images are accessible

## Maintenance

### Updating the MCP Server

1. Modify code in `dfc/mcp-server/`
2. Push to main or create a PR
3. Once merged, the container updates automatically
4. Users get updates on next plugin invocation

### Monitoring

- **GitHub Actions**: Check workflow runs in dfc repo
- **Container Registry**: View packages at github.com/chainguard-dev?tab=packages
- **Usage**: Monitor pull statistics

## Benefits of This Approach

1. **No language dependencies**: Users don't need Go, Python, Node, etc.
2. **Consistent**: Both plugins use the same pattern
3. **Secure**: Chainguard base images, signed, with SBOMs
4. **Automatic updates**: Containers rebuild on code changes
5. **Portable**: Works on any platform
6. **Versioned**: Can pin specific versions if needed
7. **Cacheable**: Layer caching speeds up pulls

## Future Enhancements

Potential improvements:

1. **Multi-registry**: Publish to both ghcr.io and cgr.dev
2. **Scheduled builds**: Rebuild weekly even without code changes
3. **Vulnerability scanning**: Integrate Trivy or Grype
4. **Performance metrics**: Track container startup time
5. **Size optimization**: Further reduce image size
6. **Release notes**: Auto-generate from commits

## References

- [GitHub Container Registry Docs](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Cosign Documentation](https://docs.sigstore.dev/cosign/overview/)
- [Chainguard Images](https://images.chainguard.dev/)
- [MCP Protocol](https://modelcontextprotocol.io/)
