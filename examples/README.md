# Examples

Five example applications that demonstrate what the `chainguard-codegen` plugin does:
convert an existing Dockerfile into a Chainguard-Containers version that is smaller,
distroless, and non-root by default.

## Convention

Every example directory ships two Dockerfiles:

- `Dockerfile.before` — the original, using Docker Hub upstream images (python-slim,
  node-bookworm-slim, golang-alpine, etc.). This is the input a user would migrate
  *from*.
- `Dockerfile` — the Chainguard-migrated version. This is what `chainguard-codegen`
  (or the `dfc` Dockerfile Converter) produces. It uses `cgr.dev/chainguard/*` base
  images, runs as UID 65532, and follows the plugin's mapping rules.

Build either to see the difference:

```bash
docker build -f examples/advisory-api/Dockerfile.before examples/advisory-api/
docker build -f examples/advisory-api/Dockerfile        examples/advisory-api/
```

## The examples

- **advisory-api** — Flask + gunicorn service returning a list of sample CVE
  advisories. Shows a multi-stage Python migration using `--target=/deps` for a
  clean runtime image. Note: `cmd.sh` is consumed by `Dockerfile.before` only;
  the Chainguard version invokes gunicorn directly via its `ENTRYPOINT`.
- **debian-toolbox** — minimal toolbox image with `curl`, `wget`, and
  `ca-certificates`. Shows `debian:bookworm-slim` → `cgr.dev/chainguard/wolfi-base`,
  `apt-get` → `apk`, and Busybox-style `addgroup`/`adduser`.
- **python-pinned** — tiny Flask app whose original Dockerfile pins a specific
  Docker Hub digest. The Chainguard version pins a digest on
  `cgr.dev/chainguard/python:latest-dev`, demonstrating how to preserve strict
  reproducibility after migration.
- **sbom-scanner** — Go HTTP service producing a static binary. Shows
  `golang:1.22-alpine` → `cgr.dev/chainguard/go:latest-dev` for the build stage
  and `scratch` → `cgr.dev/chainguard/static:latest` for the runtime. The
  `static` image already ships CA certificates and runs as UID 65532, so the
  migrated Dockerfile is noticeably simpler.
- **vuln-dashboard** — small Node.js Express app. Shows
  `node:20.18.0-bookworm-slim` → `cgr.dev/chainguard/node:latest-dev` for the
  build stage and `cgr.dev/chainguard/node:latest` for the runtime.

## See also

- The `chainguard-codegen` plugin: `../chainguard-codegen/`
- The upstream dfc tool: https://github.com/chainguard-dev/dfc
- Chainguard Containers catalog: https://images.chainguard.dev

## A note on image tags

The migrated Dockerfiles use `:latest` and `:latest-dev`, which is the
recommended Chainguard default — Chainguard rebuilds its images continuously and
the `:latest*` tags always point at the freshest CVE-patched build. Version
tags like `:3.11` and `:3.11-dev` exist on paid-tier Chainguard catalogs; the
free/public tier carries `:latest*` only.

For stricter reproducibility, pin a digest after `:latest` or `:latest-dev` as
`python-pinned/Dockerfile` does. Refresh the digest periodically since older
ones age out of the registry's retention window.
