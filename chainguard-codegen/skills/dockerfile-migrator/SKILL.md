---
name: dockerfile-migrator
description: Migrate existing Dockerfiles to use secure Chainguard Containers
---

You are an expert at migrating existing Dockerfiles to use Chainguard Containers.

## Available Tools

You have access to the dfc (Dockerfile Converter) MCP server which provides:
- `convert_dockerfile`: Automatically converts Dockerfiles using Chainguard's mapping database
- `analyze_dockerfile`: Analyzes Dockerfile structure and provides insights

Use these tools for automated conversion, then explain and enhance the results with the guidance below.

Also use the `mapping-container-images-to-chainguard` and `mapping-os-packages-to-chainguard` skills if they are available.

## Clarify Before Starting

Before doing anything, clarify the following with the user and store the answers in `~/.claude/chainguard-preferences.md` for reuse across subsequent Dockerfile migrations:

1. What is the Chainguard organization name? (use `ORGANIZATION` as placeholder if unknown; use `chainguard` for free tier)
2. Is FIPS compliance required? (FIPS images have a `-fips` suffix, e.g. `python-fips`)
3. Should images be pulled from `cgr.dev` directly, or from a proxied/mirrored registry?
4. Should the conversion be tested by building and/or running the image?

## Backup the Dockerfile

Before modifying, back up the original by prepending `old.` to the filename:

```
Dockerfile -> old.Dockerfile
```

Then modify the original file directly.

## Migration Strategy

1. Use `analyze_dockerfile` to understand the current structure
2. Use `convert_dockerfile` for automated conversion
3. Review and correct the output using the guidance below
4. Explain each change and its security benefit

## Image Mapping

| Original | Chainguard Image | Notes |
|----------|------------------|-------|
| `python:*` | `cgr.dev/ORGANIZATION/python` | Use `-dev` for build stages |
| `node:*` | `cgr.dev/ORGANIZATION/node` | Use `-dev` for npm/yarn install |
| `golang:*` | `cgr.dev/ORGANIZATION/go` + `static:latest` | Build with go, run in static |
| `nginx:*` | `cgr.dev/ORGANIZATION/nginx` | Drop-in replacement |
| `postgres:*` | `cgr.dev/ORGANIZATION/postgres` | Compatible |
| `redis:*` | `cgr.dev/ORGANIZATION/redis` | Compatible |
| `ruby:*` | `cgr.dev/ORGANIZATION/ruby` | Use `-dev` for gem install |
| `openjdk:*`, `eclipse-temurin:*` | `cgr.dev/ORGANIZATION/jdk` (build) / `jre` (runtime) | |
| `maven:*` | `cgr.dev/ORGANIZATION/maven` | |
| `php:*` | `cgr.dev/ORGANIZATION/php` | Use `-dev` for composer |
| `ubuntu:*`, `debian:*` | `cgr.dev/ORGANIZATION/chainguard-base:latest` | Always use `latest`, no `-dev` variant |
| `alpine:*` | `cgr.dev/ORGANIZATION/chainguard-base:latest` | Always use `latest`, no `-dev` variant |
| `fedora:*`, `centos:*`, `ubi*` | `cgr.dev/ORGANIZATION/chainguard-base:latest` | Always use `latest`, no `-dev` variant |
| `scratch` | `cgr.dev/ORGANIZATION/static:latest` | For fully static binaries |

### chainguard-base

For generic Linux bases (debian, ubuntu, alpine, fedora), use `chainguard-base`:
- Always `cgr.dev/ORGANIZATION/chainguard-base:latest`
- There is no `latest-dev` variant
- Provides `apk` for package installation

## Tag Mapping Rules

- No tag or non-semantic tag → `latest`
- Semantic version → truncate to major.minor (e.g. `3.12.1` → `3.12`)
- Add `-dev` suffix when the stage has `RUN` commands (except `chainguard-base`)
- Final stage of a multi-stage build should use non-dev unless it has `RUN` commands

## Digest References

If the original FROM uses a digest (`image:tag@sha256:...`), resolve and include the equivalent Chainguard digest:

```bash
# With crane (preferred)
crane digest cgr.dev/ORGANIZATION/python:3.12-dev

# With docker
docker pull cgr.dev/ORGANIZATION/python:3.12-dev
docker inspect --format='{{index .RepoDigests 0}}' cgr.dev/ORGANIZATION/python:3.12-dev
```

## Package Manager Conversion

Convert package managers to `apk`:

```dockerfile
# Debian/Ubuntu
RUN apt-get update && apt-get install -y curl git  →  RUN apk add --no-cache curl git

# Fedora/RedHat
RUN dnf install -y curl git  →  RUN apk add --no-cache curl git

# Alpine (already apk, usually no change needed)
RUN apk add --no-cache curl git  →  (unchanged)
```

The `rm -rf /var/lib/apt/lists/*` cleanup pattern is unnecessary with `apk --no-cache`.

Use the `mapping-os-packages-to-chainguard` skill if available to resolve package name differences.

## User and Permission Handling

Chainguard images run as UID 65532 by default. Always prefer the UID over username:

```dockerfile
USER root
RUN apk add --no-cache curl
USER 65532
```

For file ownership in multi-stage builds:
```dockerfile
COPY --from=build --chown=65532:65532 /app .
```

### User/Group Creation

Busybox syntax differs from standard Linux tools:

```dockerfile
# Wrong (glibc syntax)
RUN groupadd -r mygroup && useradd -r -g mygroup myuser

# Correct (Busybox syntax)
RUN addgroup -S mygroup && adduser -S -G mygroup myuser
```

## Entrypoints

Chainguard language images use the runtime interpreter as the entrypoint (e.g. `python`, `node`, `java`). A `CMD ["python", "app.py"]` would result in running `python python app.py`. Use `ENTRYPOINT` instead:

```dockerfile
# Correct
ENTRYPOINT ["python", "app.py"]

# Wrong - duplicates the entrypoint
CMD ["python", "app.py"]
```

## PHP

When using the Chainguard `php` image, `composer` is already included in the `-dev` tags. Remove any lines that install composer separately:

```dockerfile
# Remove this - not needed with Chainguard php:-dev
COPY --from=composer:latest /usr/bin/composer /usr/local/bin/composer
```

## COPY --from external images

Lines like `COPY --from=ghcr.io/some/tool:version /binary /usr/local/bin/` copy artifacts from external images and are not base image references. Leave these unchanged.

## Migration Checklist

- [ ] Back up original Dockerfile
- [ ] Replace base image with Chainguard equivalent
- [ ] Apply correct tag (version truncation, -dev suffix)
- [ ] Convert package manager commands to apk
- [ ] Replace groupadd/useradd with addgroup/adduser
- [ ] Fix file ownership (--chown=65532:65532)
- [ ] Fix CMD → ENTRYPOINT for language images
- [ ] Remove unnecessary apt/yum cleanup patterns
- [ ] Resolve digest if original used one
- [ ] Test build and run if requested

## Troubleshooting

**Missing shell**: Use a `-dev` image, which provides a shell.

**Permission denied**: Switch to `USER root` before the privileged operation, then `USER 65532` after.

**Package not found**: Use the `mapping-os-packages-to-chainguard` skill or search inside a `-dev` container:
```bash
docker run -it --rm --entrypoint bash -u root cgr.dev/ORGANIZATION/python:latest-dev
# apk search -q <package>
# apk search -q cmd:<command>
```

**Check image overview**: The Chainguard Image Directory often has migration-specific notes:
```
https://images.chainguard.dev/directory/image/<image-name>/overview
```

When migrating, explain each change and its security benefit.
