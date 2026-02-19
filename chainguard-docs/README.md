# Chainguard Documentation Plugin

Access Chainguard's complete documentation library directly within Claude Code through a Model Context Protocol (MCP) server.

## Features

- **Search Documentation**: Find relevant content across all Chainguard docs
- **Container Images**: Get detailed information about specific Chainguard Containers
- **Security Guides**: Access CVE management, SBOM, and signing documentation
- **Tool References**: Learn about Wolfi, apko, melange, and chainctl
- **Real-time Updates**: Documentation is automatically updated weekly

## Installation

Users can install this plugin by adding the Chainguard marketplace:

```bash
/plugin marketplace add github.com/chainguard-demo/claude-plugins
/plugin install chainguard-docs@chainguard-plugins
```

## Requirements

- Docker installed and running
- Internet connection to pull the MCP server image

## Available MCP Tools

Once installed, Claude Code can use these tools automatically:

### `search_docs`
Search across all Chainguard documentation for relevant content.

**Parameters:**
- `query` (string): Search query
- `max_results` (integer, optional): Maximum results (default: 5)

### `get_image_docs`
Get documentation for a specific Chainguard container image.

**Parameters:**
- `image_name` (string): Image name (e.g., "python", "node", "nginx")

### `list_images`
List all available Chainguard container images.

### `get_security_docs`
Get security-related documentation including CVE management and SBOMs.

### `get_tool_docs`
Get documentation for Chainguard tools.

**Parameters:**
- `tool_name` (string): Tool name (wolfi, apko, melange, or chainctl)

## Example Usage

Once installed, you can ask Claude Code:

```
"What Chainguard images are available for Python?"
"Show me how to use apko"
"Search Chainguard docs for SBOM generation"
"How does Chainguard handle CVEs?"
```

Claude Code will automatically use the MCP tools to retrieve relevant documentation.

## Documentation Source

The MCP server provides access to:
- Container image guides from [images.chainguard.dev](https://images.chainguard.dev)
- Security and tooling documentation from [edu.chainguard.dev](https://edu.chainguard.dev)
- Best practices and tutorials

## Support

- [Chainguard Academy](https://edu.chainguard.dev)
- [Chainguard Support](https://support.chainguard.dev)
- [Community Slack](https://go.chainguard.dev/slack)
