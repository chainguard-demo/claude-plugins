# Build stage: compile a static binary using the Chainguard Go image.
FROM cgr.dev/chainguard/go:latest-dev AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 is required for a cgr.dev/chainguard/static runtime — no libc.
RUN CGO_ENABLED=0 GOOS=linux go build -o sbom-scanner .

# Runtime stage: distroless static image, runs as nonroot (UID 65532) by default.
FROM cgr.dev/chainguard/static:latest

COPY --from=builder /app/sbom-scanner /sbom-scanner

EXPOSE 8080

ENTRYPOINT ["/sbom-scanner"]
