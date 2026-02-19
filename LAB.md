# Migrate to Chainguard with Claude — Lab Guide

This lab walks through migrating container workloads to Chainguard Containers using
Claude Code plugins, skills, and Chainguard Academy's AI-optimized documentation.

**Prerequisites:**
- Docker installed and running
- Claude Code installed (`npm install -g @anthropic-ai/claude-code`)
- `cosign` installed (optional, for bundle verification)

---

## Part 1 — Getting Chainguard Docs into AI Tools

Before we touch any tooling, I want to show you something
that you can use right now, with any AI tool — not just Claude. Chainguard
Academy has been built from the ground up to work well with AI assistants.
There are several ways to get the docs into your AI tool of choice, ranging
from a single click to a fully integrated MCP server. Let's walk through them.

---

### 1a — Per-page Markdown copy

The simplest option. Every page on Chainguard Academy has a
'Copy LLM text' button that gives you the entire page as clean Markdown — no
HTML noise, no navigation chrome. You paste it into Claude, ChatGPT, Gemini,
whatever you're using. This is great when you have a specific question about
one image or one topic.

**Demo steps:**
1. Open `https://edu.chainguard.dev/chainguard/migration/migrating-to-chainguard-images/`
2. Point out the **Copy LLM text** button
3. Click it, paste into Claude.ai, ask: *"Based on this documentation, what's the
   first thing I should do before migrating a Dockerfile?"*

---

### 1b — `llms.txt` and `llms-full.txt`

If you want broader context, we publish two machine-readable
files that follow the emerging `llms.txt` standard. The lightweight version is
a structured index — good for helping an AI navigate the site. The full version
is the entire documentation corpus compiled at build time. It's large, but if
you're doing a deep migration project it means you can paste everything into a
long-context model and ask questions across the whole docs set.

```bash
# Lightweight index — site structure and summaries
curl https://edu.chainguard.dev/llms.txt

# Full documentation corpus — everything, compiled at build time
curl -s https://edu.chainguard.dev/llms-full.txt | wc -c   # check the size
curl -O https://edu.chainguard.dev/llms-full.txt
```

You can drop that file into any AI tool that supports file
uploads, or reference it from a system prompt in your own AI pipelines.

---

### 1c — Static extraction and the signed docs bundle

For teams building internal tooling, RAG pipelines, or
anything where you need the docs as a local artifact, there are two options.
You can extract directly from the same container we use for the MCP server, or
you can download the signed weekly release from GitHub. I'll show you both —
and the signed release is worth calling out specifically because it comes with
Cosign signatures and SLSA Level 3 provenance. You can verify exactly what
you're feeding to your AI.

**Option A — Extract from the container:**
```bash
docker pull ghcr.io/chainguard-dev/ai-docs:latest

# Extracts all documentation as a single markdown file into ./output/
docker run --rm -v $(pwd):/output \
  ghcr.io/chainguard-dev/ai-docs:latest extract /output

ls -lh output/
```

**Option B — Download the signed weekly release:**
```bash
# Download the bundle and its signature/certificate
curl -LO https://github.com/chainguard-dev/edu/releases/download/ai-docs-latest/chainguard-ai-docs.tar.gz
curl -LO https://github.com/chainguard-dev/edu/releases/download/ai-docs-latest/chainguard-ai-docs.tar.gz.sig
curl -LO https://github.com/chainguard-dev/edu/releases/download/ai-docs-latest/chainguard-ai-docs.tar.gz.crt

# Verify the signature with cosign
cosign verify-blob \
  --certificate chainguard-ai-docs.tar.gz.crt \
  --signature chainguard-ai-docs.tar.gz.sig \
  --certificate-identity "https://github.com/chainguard-dev/edu/.github/workflows/compile-public-docs.yml@refs/heads/main" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  chainguard-ai-docs.tar.gz

# Extract
tar -xzf chainguard-ai-docs.tar.gz
ls -lh
```

This is Chainguard being Chainguard — the documentation you
feed to your AI has provenance. You know it came from our build pipeline, it
hasn't been tampered with, and you can verify that independently. More at
`edu.chainguard.dev/ai-docs-security`.

---

### 1d — Clone the Academy repo and work directly in it

Here's something most people don't realise. Chainguard Academy
is a public GitHub repo — all the documentation is just Markdown files. So
another way to get the docs into an AI tool is to clone the repo and open
Claude Code directly inside it. Claude then has native file access to the
entire content tree — no container, no curl, no paste. It can read, search,
and cross-reference any page instantly. This is particularly useful if you're
a documentation contributor, or if you want to ask questions that span multiple
areas of the docs at once.

```bash
git clone https://github.com/chainguard-dev/edu
cd edu
claude
```

Once you're inside, you can ask questions across the full
corpus just by referencing the content. For example:

```
What does the migration guide say about multi-stage builds?
```

```
Summarise the differences between the chainguard-base and wolfi-base images
based on the documentation in this repo.
```

```
Find all the pages that mention FIPS and give me a one-line summary of each.
```

The key advantage here over pasting llms-full.txt into a chat
is that Claude Code can navigate the file tree, follow cross-references, and
work incrementally — it doesn't need to load everything into context at once.
It also means any edits you make to the docs are immediately visible to Claude,
which is great for doc contributors reviewing their own changes.

---

### 1e — MCP server (preview of where we're going)

The same container has one more mode. Instead of extracting
static files, it can run as an MCP server — which is what lets AI assistants
query the docs dynamically, only retrieving what's relevant to the current
question. This is what the Claude Code plugin wraps. Let's set that up now.

```bash
# Test it directly — you'll see the server initialize then wait on stdin
docker run --rm -i ghcr.io/chainguard-dev/ai-docs:latest serve-mcp
# Ctrl+C to exit
```

> **Apple Silicon (M1/M2/M3/M4):** The image is currently amd64-only. Add `--platform linux/amd64` to run it via Rosetta emulation:
> ```bash
> docker run --rm -i --platform linux/amd64 ghcr.io/chainguard-dev/ai-docs:latest serve-mcp
> ```

---

## Part 2 — Claude Code Plugin Setup

Now let's install the Chainguard plugins for Claude Code.
There are two: `chainguard-docs` wraps that MCP server we just saw, and
`chainguard-codegen` adds the dfc Dockerfile Converter as an MCP server plus
a set of migration skills. Once these are installed, Claude has live access
to Chainguard documentation and can convert Dockerfiles automatically.

```bash
# Add the Chainguard plugin marketplace
/plugin marketplace add github.com/chainguard-demo/claude-plugins

# Install both plugins
/plugin install chainguard-docs@chainguard-plugins
/plugin install chainguard-codegen@chainguard-plugins
```

**Demo — chainguard-docs in action:**

Ask Claude Code:
```
What Chainguard images are available for Python, and what's the difference
between :latest and :latest-dev?
```

```
How does Chainguard handle CVE management?
```

```
What's the Chainguard image equivalent of nginx:alpine?
```

Notice that Claude is calling `get_image_docs` and
`search_docs` against the live MCP server rather than relying on its training
data. That matters because Chainguard Containers are updated daily — you want
current information.

---

## Part 3 — Simple Migration: advisory-api

Let's migrate a real Dockerfile. We'll start with something
small — advisory-api, a single-stage Python Flask service for serving security
advisories. It uses Python 3.11, creates a user with `groupadd` and `useradd`,
and installs packages with `pip` inline. Simple enough that we can understand
every change the migration makes.

```bash
# Take a look at the original
cat examples/advisory-api/Dockerfile
```

Expected output:
```dockerfile
FROM python:3.11-slim

RUN groupadd -r advisory && useradd -r -g advisory advisory
RUN pip install flask==3.0.0 gunicorn==21.2.0 requests==2.31.0 redis==5.0.0
WORKDIR /app
COPY app /app
COPY cmd.sh /

EXPOSE 9090
USER advisory

CMD ["/cmd.sh"]
```

Now ask Claude Code (with the plugin installed):
```
Migrate examples/advisory-api/Dockerfile to use Chainguard Containers.
My organization is chainguard.
```

Walk through what changed:
- `python:3.11-slim` → `cgr.dev/chainguard/python:3.11-dev` — version preserved,
  `-dev` added because the stage has RUN commands
- `groupadd`/`useradd` → `addgroup`/`adduser` — Busybox syntax in Wolfi
- `pip install` stays but gets `USER root` before it and `USER 65532` after
- `CMD ['/cmd.sh']` is fine here since it's a shell script, not a language runtime

Notice Claude backed up the original as `old.Dockerfile` first — that's the
skill telling it to do that before making changes.

---

## Part 4 — Complex Migration: vuln-dashboard

Now something more representative of a real production workload.
vuln-dashboard is a two-stage Node.js build — a builder stage that compiles dependencies and a minimal runtime stage. It has apt-get in both stages, `groupadd`/`useradd` with ARG-based UIDs, and enough structure to show where AI adds real value over running dfc manually.

```bash
cat examples/vuln-dashboard/Dockerfile
```

Ask Claude Code:
```
Migrate examples/vuln-dashboard/Dockerfile to use Chainguard Containers.
My organization is chainguard. Please test by building the image after migration.
```

Things to highlight as Claude works through it:

- **dfc via MCP handles the bulk** — FROM conversions, apt-get → apk, USER root wrapping
- **Version truncation** — `node:20.18.0-bookworm-slim` → `node:20.18-dev`
- **Both stages get `-dev`** — both stages have RUN commands
- **groupadd/useradd → addgroup/adduser** — converted in both stages
- **apt cleanup removed** — `rm -rf /var/lib/apt/lists/*` is unnecessary with `apk --no-cache`

---

## Part 5 — Skills: Deeper Migration Guidance (optional / advanced)

The plugins include three standalone skills that give Claude
more detailed guidance for edge cases — things like FIPS compliance, digest
pinning, mapping packages that don't have obvious equivalents, and knowing
when to use `chainguard-base` vs a specific language image. These work with
Claude Code independently of the plugin, and because they're SKILL.md format
they also work with Cursor, Gemini CLI, and other compatible tools.

```bash
# Install the standalone skills
cp -r chainguard-codegen/skills/migrating-dockerfiles-to-chainguard ~/.claude/skills/
cp -r chainguard-codegen/skills/mapping-container-images-to-chainguard ~/.claude/skills/
cp -r chainguard-codegen/skills/mapping-os-packages-to-chainguard ~/.claude/skills/
```

**Demo — package mapping:**
```bash
cat examples/debian-toolbox/Dockerfile
```

Ask Claude Code:
```
Migrate examples/debian-toolbox/Dockerfile to Chainguard.
The Dockerfile installs curl, wget, and ca-certificates — verify the package
names are correct for Wolfi before migrating.
```

The `mapping-os-packages-to-chainguard` skill will kick in
here. It downloads the dfc builtin mappings YAML and uses yq to look up each
package name. Most are the same in Wolfi, but for anything non-obvious it
will search inside a running `-dev` container with `apk search`.

**Demo — digest pinning:**
```bash
cat examples/python-pinned/Dockerfile
```

Ask Claude Code:
```
Migrate examples/python-pinned/Dockerfile to Chainguard.
Preserve the digest pin.
```

The `mapping-container-images-to-chainguard` skill handles
this — it knows to resolve a new digest for the Chainguard image using crane
or docker inspect, rather than leaving a stale hash or dropping the pin entirely.

---

## Reference

| Resource | URL / Command |
|----------|---------------|
| Chainguard Academy | `https://edu.chainguard.dev` |
| Developer resources | `https://edu.chainguard.dev/developer-resources` |
| AI docs security | `https://edu.chainguard.dev/ai-docs-security` |
| llms.txt | `https://edu.chainguard.dev/llms.txt` |
| llms-full.txt | `https://edu.chainguard.dev/llms-full.txt` |
| Signed docs bundle | `https://github.com/chainguard-dev/edu/releases/tag/ai-docs-latest` |
| Plugin marketplace | `github.com/chainguard-demo/claude-plugins` |
| ai-docs container | `ghcr.io/chainguard-dev/ai-docs:latest` |
| dfc-mcp container | `cgr.dev/chainguard/dfc-mcp` |
| Image directory | `https://images.chainguard.dev` |
| Migration checklist PDF | `https://edu.chainguard.dev/downloads/migrating-to-chainguard-images.pdf` |
