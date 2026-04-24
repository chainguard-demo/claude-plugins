FROM cgr.dev/chainguard/wolfi-base:latest

USER root
RUN apk add --no-cache curl wget ca-certificates

WORKDIR /app
COPY --chown=65532:65532 . .
RUN chmod +x entrypoint.sh

USER 65532

ENTRYPOINT ["./entrypoint.sh"]
