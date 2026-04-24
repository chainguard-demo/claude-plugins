# Build stage: install dependencies into an isolated directory
FROM cgr.dev/chainguard/python:latest-dev AS builder

USER root
RUN pip install --no-cache-dir --target=/deps \
    flask==3.0.0 \
    gunicorn==21.2.0 \
    requests==2.31.0 \
    redis==5.0.0

# Runtime stage: minimal distroless image, no shell required
FROM cgr.dev/chainguard/python:latest

COPY --from=builder --chown=65532:65532 /deps /deps
COPY --chown=65532:65532 app /app

ENV PYTHONPATH=/deps

WORKDIR /app
EXPOSE 9090

# Invoke gunicorn directly from /deps/bin so the runtime doesn't rely on
# `python -m gunicorn` indirection. Overrides the base image's python entrypoint.
ENTRYPOINT ["/deps/bin/gunicorn", "--bind", "0.0.0.0:9090", "main:app"]
