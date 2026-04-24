# Build stage: install dependencies
FROM cgr.dev/chainguard/node:latest-dev AS builder

WORKDIR /app

USER root
COPY package*.json ./
RUN npm install
COPY . .

# Runtime stage: minimal distroless image
FROM cgr.dev/chainguard/node:latest

WORKDIR /app

COPY --from=builder --chown=65532:65532 /app/node_modules ./node_modules
COPY --from=builder --chown=65532:65532 /app/server.js ./

EXPOSE 3000

# Chainguard node image entrypoint is `node`, CMD args are passed directly to it
CMD ["server.js"]
