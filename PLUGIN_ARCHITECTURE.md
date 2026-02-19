# Chainguard Claude Code Plugins - Architecture

This document explains the architecture and design of the Chainguard plugins for Claude Code.

## Overview

The Chainguard plugin marketplace consists of two complementary plugins:

1. **chainguard-docs**: Documentation access via MCP
2. **chainguard-codegen**: Code generation with dfc MCP integration

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Claude Code                              │
│                                                             │
│  ┌───────────────────────────┐  ┌─────────────────────────┐│
│  │  chainguard-docs          │  │  chainguard-codegen     ││
│  │  ┌─────────────────────┐  │  │  ┌───────────────────┐ ││
│  │  │ MCP Client          │  │  │  │ MCP Client        │ ││
│  │  └──────────┬──────────┘  │  │  └─────────┬─────────┘ ││
│  └─────────────┼─────────────┘  └────────────┼───────────┘│
│                │                              │             │
└────────────────┼──────────────────────────────┼─────────────┘
                 │                              │
                 │ stdio                        │ stdio
                 ▼                              ▼
    ┌────────────────────────┐    ┌────────────────────────┐
    │ Chainguard AI Docs     │    │ dfc MCP Server         │
    │ MCP Server             │    │                        │
    │ (Docker container)     │    │ (go run)               │
    │                        │    │                        │
    │ Tools:                 │    │ Tools:                 │
    │ - search_docs          │    │ - convert_dockerfile   │
    │ - get_image_docs       │    │ - analyze_dockerfile   │
    │ - list_images          │    │ - healthcheck          │
    │ - get_security_docs    │    │                        │
    │ - get_tool_docs        │    │ Uses:                  │
    │                        │    │ - builtin-mappings.yaml│
    │ Serves:                │    │ - dfc Go library       │
    │ - chainguard-ai-docs.md│    │                        │
    └────────────────────────┘    └────────────────────────┘
```

## Plugin 1: chainguard-docs

### Purpose
Provides access to Chainguard's complete documentation library through an MCP server.

### Components

**MCP Server**: `ghcr.io/chainguard-dev/ai-docs:latest`
- Pre-built Docker container
- Contains indexed documentation
- Updates weekly automatically

**Files**:
```
chainguard-docs/
├── .claude-plugin/
│   └── plugin.json          # Plugin metadata
├── .mcp.json                # MCP server configuration
└── README.md                # User documentation
```

### How It Works

1. User installs plugin via Claude Code
2. Claude Code starts the MCP server (Docker container) on demand
3. When user asks documentation questions, Claude Code:
   - Calls MCP server tools (search_docs, get_image_docs, etc.)
   - Receives relevant documentation
   - Synthesizes answer for the user

### MCP Tools Provided

- `search_docs(query, max_results)` - Full-text search
- `get_image_docs(image_name)` - Specific image documentation
- `list_images()` - Available Chainguard Containers
- `get_security_docs()` - CVE, SBOM, signing info
- `get_tool_docs(tool_name)` - Tool documentation (wolfi, apko, melange, chainctl)

## Plugin 2: chainguard-codegen

### Purpose
Generates secure code configurations using Chainguard Containers and the dfc (Dockerfile Converter) tool.

### Components

**MCP Server**: dfc MCP server
- Runs via `cgr.dev/chainguard/dfc-mcp` Docker container
- Uses Chainguard's official image and package mappings
- Performs intelligent Dockerfile conversion

**AI Skills**:
- `dockerfile-generator.json` - Guides Dockerfile creation
- `dockerfile-migrator.json` - Guides Dockerfile migration (uses dfc MCP)
- `apko-generator.json` - Guides apko config creation

**Files**:
```
chainguard-codegen/
├── .claude-plugin/
│   └── plugin.json          # Plugin metadata
├── .mcp.json                # dfc MCP server configuration
├── README.md                # User documentation
├── agents/                  # (empty - ready for future agents)
├── commands/                # (empty - ready for slash commands)
└── skills/
    ├── dockerfile-generator.json
    ├── dockerfile-migrator.json
    └── apko-generator.json
```

### How It Works

#### Dockerfile Migration Flow

1. User: "Migrate this Dockerfile to use Chainguard Containers"
2. Claude Code:
   - Activates `dockerfile-migrator` skill (provides context/guidance)
   - Calls dfc MCP `convert_dockerfile` tool with Dockerfile content
   - dfc performs conversion using built-in mappings
   - Claude explains the changes and security benefits

#### Dockerfile Generation Flow

1. User: "Generate a Dockerfile for a Python Flask app"
2. Claude Code:
   - Activates `dockerfile-generator` skill
   - Generates Dockerfile following skill instructions:
     - Multi-stage build
     - Chainguard Containers
     - Non-root user
     - Security best practices

### MCP Tools Provided

- `convert_dockerfile(dockerfile_content, organization, registry)` - Full conversion
- `analyze_dockerfile(dockerfile_content)` - Structure analysis
- `healthcheck()` - Server health check

## Skills vs MCP Tools

### Skills (Passive Knowledge)
- JSON files with instructions
- Loaded by Claude Code automatically
- Guide Claude's behavior and responses
- No code execution

### MCP Tools (Active Capabilities)
- Executable functions
- Called by Claude Code via MCP protocol
- Perform actual operations
- Return structured data

### Combined Approach

The chainguard-codegen plugin uses **both**:

1. **MCP tools** (dfc) provide the conversion engine
2. **Skills** provide the expertise layer:
   - Best practices
   - Security guidance
   - Pattern examples
   - Explanation templates

This combination allows Claude to:
- Use dfc for automated conversion
- Enhance results with expert knowledge
- Explain security benefits
- Handle edge cases the automated tool might miss

## User Experience

### Installing

```bash
# Add marketplace
/plugin marketplace add github.com/chainguard-demo/claude-plugins

# Install both plugins
/plugin install chainguard-docs@chainguard-plugins
/plugin install chainguard-codegen@chainguard-plugins
```

### Using chainguard-docs

Natural language queries:
```
"What Chainguard images are available for Node.js?"
"How does Chainguard handle CVE management?"
"Show me the nginx image documentation"
```

Claude Code automatically:
1. Recognizes it needs documentation
2. Calls the appropriate MCP tool
3. Presents the information

### Using chainguard-codegen

Natural language requests:
```
"Generate a Dockerfile for a Django application"
"Migrate this Dockerfile to use Chainguard Containers"
"Create an apko config for a custom Python image"
```

Claude Code automatically:
1. Activates relevant skills
2. Calls dfc MCP tools when appropriate
3. Generates or converts code
4. Explains the security benefits

## Extension Points

Both plugins are designed for future expansion:

### chainguard-docs
- Additional MCP tools as documentation grows
- Integration with Chainguard Enforce API
- Real-time CVE lookup

### chainguard-codegen
- **Slash commands**: `/chainguard-dockerfile`, `/chainguard-migrate`
- **Agents**: Automated scanning and migration workflows
- **Additional skills**: melange configs, chainctl commands
- **MCP tools**: Direct API integration for image validation

## Dependencies

### chainguard-docs
- Docker (to run MCP server container)
- Internet access (to pull ghcr.io/chainguard-dev/ai-docs)

### chainguard-codegen
- Go 1.20+ (to run dfc MCP server via `go run`)
- Internet access (to download dfc MCP server)

Alternative: Users can build and configure local paths for both MCP servers.

## Security Considerations

### Sandboxing
- MCP servers run in isolated processes
- No direct access to user's filesystem
- Communication only via stdio protocol

### Trust
- chainguard-docs: Official Chainguard container from ghcr.io
- chainguard-codegen: Official dfc tool from Chainguard GitHub

### Data Privacy
- Both plugins operate on data provided in the conversation
- No telemetry or external API calls (except for MCP server itself)
- Documentation is included in the container (no runtime fetching)

## Testing Locally

Before publishing, test locally:

```bash
cd chainguard-claude-plugins

# Test in Claude Code
/plugin marketplace add .
/plugin install chainguard-docs@chainguard-plugins
/plugin install chainguard-codegen@chainguard-plugins

# Test queries
"List Chainguard images"
"Convert this Dockerfile to use Chainguard..."
```

## Distribution

### Community Marketplace (Current)
1. Push to GitHub: `github.com/chainguard-demo/claude-plugins`
2. Users add: `/plugin marketplace add github.com/chainguard-demo/claude-plugins`
3. Full control over updates and versioning

### Official Directory (Future)
1. Submit to anthropics/claude-plugins-official
2. PR to `/external_plugins/chainguard-*`
3. Requires Anthropic approval
4. Listed in official plugin directory

## Maintenance

### chainguard-docs
- MCP server updated weekly (automatic)
- Plugin configuration rarely changes
- Documentation stays current automatically

### chainguard-codegen
- dfc MCP server updated with dfc releases
- Skills may need updates for new features
- Image mappings updated upstream in dfc

## Future Enhancements

1. **Slash Commands**
   - `/chainguard-dockerfile <stack>` - Quick generation
   - `/chainguard-migrate` - Interactive migration wizard

2. **Agents**
   - Automated repository scanning
   - Bulk Dockerfile migration
   - CI/CD integration guidance

3. **Additional Skills**
   - melange package building
   - chainctl command generation
   - Chainguard Enforce configuration

4. **Enhanced MCP Integration**
   - Direct Chainguard API access
   - Real-time image vulnerability checks
   - SBOM generation and validation
