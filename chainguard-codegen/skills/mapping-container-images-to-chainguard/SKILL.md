---
name: mapping-container-images-to-chainguard
description: Instructions for mapping container image references in FROM statements to their Chainguard equivalents. Use this when converting Dockerfiles to Chainguard.
---

When converting `FROM` statements in a Dockerfile to use Chainguard images, follow
these steps:

1. Use the dfc MCP server for automated conversion where available
2. Apply tag mapping rules
3. Look up unknown images using the ai-docs MCP server
4. Handle special cases (generic Linux bases, static binaries, FIPS)

## Use the dfc MCP Server

If the `dfc` MCP server is available, use it to perform automated conversion. It
handles registry rewriting, image name translation, and tag mapping automatically.

Use the `convert_dockerfile` tool to convert a full Dockerfile, or
`analyze_dockerfile` to inspect the current state before converting.

When using dfc, always pass the organization name if you know it. If you don't,
use `ORG` as a placeholder and remind the user to substitute their own.

## Organization Name

Every Chainguard customer has their own organization namespace at `cgr.dev`:

```
cgr.dev/ORGANIZATION/image-name:tag
```

- For enterprise customers, use the customer's Chainguard organization name
- For free tier / public images, use `chainguard`

Check `/tmp/chainguard-preferences.md` if it exists, or ask the user if the
organization name is unknown.

## Tag Mapping Rules

Follow these rules when converting tags from upstream images to Chainguard:

1. **`chainguard-base` always uses `latest`** — never add a version or `-dev` suffix
2. **No tag or non-semantic tag → `latest`**
   - `FROM node` → `FROM cgr.dev/ORGANIZATION/node:latest`
   - `FROM node:lts` → `FROM cgr.dev/ORGANIZATION/node:latest`
3. **Semantic version → truncate to major.minor**
   - `FROM node:14.17.3` → `FROM cgr.dev/ORGANIZATION/node:14.17`
4. **Add `-dev` suffix when the stage has `RUN` commands** (except `chainguard-base`)
   - `FROM node:14` (stage has RUN) → `FROM cgr.dev/ORGANIZATION/node:14-dev`
   - `FROM node:14` (no RUN) → `FROM cgr.dev/ORGANIZATION/node:14`
5. **Final stage of a multi-stage build should use a non-dev image** unless it
   has `RUN` commands

Examples:

- `FROM node:14` → `FROM cgr.dev/ORG/node:14-dev` (if stage has RUN commands)
- `FROM node:14.17.3` → `FROM cgr.dev/ORG/node:14.17-dev` (if stage has RUN commands)
- `FROM debian:bullseye` → `FROM cgr.dev/ORG/chainguard-base:latest` (always)
- `FROM golang:1.19-alpine` → `FROM cgr.dev/ORG/go:1.19-dev` (if stage has RUN commands)

## Image Mapping Table

| Upstream Image | Chainguard Image | Notes |
|---|---|---|
| `node`, `node:*` | `cgr.dev/ORGANIZATION/node` | |
| `python:*` | `cgr.dev/ORGANIZATION/python` | |
| `golang:*`, `go:*` | `cgr.dev/ORGANIZATION/go` | Use `static` for distroless runtime |
| `golang:*-alpine` | `cgr.dev/ORGANIZATION/go` | Strip the `-alpine` suffix |
| `openjdk:*`, `eclipse-temurin:*`, `amazoncorretto:*` | `cgr.dev/ORGANIZATION/jdk` (build) or `cgr.dev/ORGANIZATION/jre` (runtime) | |
| `maven:*` | `cgr.dev/ORGANIZATION/maven` | |
| `gradle:*` | `cgr.dev/ORGANIZATION/gradle` | |
| `php:*-cli`, `php:*` | `cgr.dev/ORGANIZATION/php` | |
| `php:*-fpm` | `cgr.dev/ORGANIZATION/php:latest-fpm` | Pass custom mapping to dfc via `--mapping` |
| `nginx:*` | `cgr.dev/ORGANIZATION/nginx` | |
| `ruby:*` | `cgr.dev/ORGANIZATION/ruby` | |
| `rust:*` | `cgr.dev/ORGANIZATION/rust` | |
| `debian:*`, `ubuntu:*` | `cgr.dev/ORGANIZATION/chainguard-base:latest` | Always `latest`, no `-dev` variant |
| `alpine:*` | `cgr.dev/ORGANIZATION/chainguard-base:latest` | Always `latest`, no `-dev` variant |
| `fedora:*`, `centos:*`, `ubi*` | `cgr.dev/ORGANIZATION/chainguard-base:latest` | Always `latest`, no `-dev` variant |
| `scratch` | `cgr.dev/ORGANIZATION/static:latest` | For fully static binaries |
| `gcr.io/distroless/static:*` | `cgr.dev/ORGANIZATION/static:latest` | |
| `gcr.io/distroless/base:*` | `cgr.dev/ORGANIZATION/chainguard-base:latest` | |
| `postgres:*` | `cgr.dev/ORGANIZATION/postgres` | |
| `redis:*` | `cgr.dev/ORGANIZATION/redis` | |
| `mysql:*`, `mariadb:*` | `cgr.dev/ORGANIZATION/mariadb` | |
| `memcached:*` | `cgr.dev/ORGANIZATION/memcached` | |

## Special Cases

### Generic Linux Bases

For generic Linux distributions (`debian`, `ubuntu`, `alpine`, `fedora`, `centos`,
`ubi*`), use `chainguard-base`:

```
FROM cgr.dev/ORGANIZATION/chainguard-base:latest
```

- Always use the `latest` tag
- There is no `latest-dev` variant for `chainguard-base`
- `chainguard-base` provides `apk` for package installation

### Static and Distroless Runtimes

For compiled binaries (Go, Rust, C) with no external runtime dependencies:

```dockerfile
# Build stage
FROM cgr.dev/ORGANIZATION/go:latest-dev AS builder
# ... build with CGO_ENABLED=0 for fully static ...

# Runtime stage
FROM cgr.dev/ORGANIZATION/static:latest
```

Note that `CGO_ENABLED=0` is typically required for the `static` image. If the
binary links against glibc, check `images.chainguard.dev` for a suitable
dynamic runtime image.

### FIPS Variants

If FIPS compliance is required (check `/tmp/chainguard-preferences.md`), append
`-fips` to the image name:

- `python` → `python-fips`
- `node` → `node-fips`
- `go` → `go-fips`
- `jdk` → `jdk-fips`
- `chainguard-base` → `chainguard-base-fips`

For example: `cgr.dev/ORGANIZATION/python-fips:latest-dev`

## Digest References

If the original `FROM` uses a digest reference, the converted image should also
include a digest. The digest will naturally be different for the Chainguard image.

Example original:
```
FROM python:3.12-slim@sha256:a866731a6b71c4a194a845d86e06568725e430ed21821d0c52e4efb385cf6c6f
```

Should become something like:
```
FROM cgr.dev/ORGANIZATION/python:3.12-dev@sha256:<chainguard-digest>
```

To get the correct digest for the Chainguard image:

```
# With crane (prefer this if available)
crane digest cgr.dev/ORGANIZATION/python:3.12-dev

# Or with docker
docker pull cgr.dev/ORGANIZATION/python:3.12-dev
docker inspect --format='{{index .RepoDigests 0}}' cgr.dev/ORGANIZATION/python:3.12-dev
```

## Looking Up Images with the ai-docs MCP Server

If you're unsure which Chainguard image to use, or need to verify an image exists
and check its documentation, use the `ai-docs` MCP server:

- `list_images()` — see all available Chainguard images
- `get_image_docs(image_name)` — get documentation for a specific image
- `search_docs(query)` — search by technology, use case, or keyword

For example, to find the right image for a Laravel PHP application:

```
get_image_docs("php")
search_docs("PHP FPM Chainguard")
```

You can also check the Chainguard Image Directory directly:

```
https://images.chainguard.dev/directory/image/<image-name>/overview
https://images.chainguard.dev/directory/image/<image-name>/versions
```
