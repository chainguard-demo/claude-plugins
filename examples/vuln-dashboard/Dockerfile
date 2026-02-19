FROM node:20.18.0-bookworm-slim AS builder

WORKDIR /app

ARG APP_UID=1001
ARG APP_GID=1001

RUN apt-get update \
  && apt-get install -y --no-install-recommends build-essential \
  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man \
  && apt-get clean \
  && groupadd -g "${APP_GID}" appuser \
  && useradd --create-home --no-log-init -u "${APP_UID}" -g "${APP_GID}" appuser \
  && chown appuser:appuser /app

USER appuser

COPY --chown=appuser:appuser package*.json ./
RUN npm install

COPY --chown=appuser:appuser . .

FROM node:20.18.0-bookworm-slim

WORKDIR /app

ARG APP_UID=1001
ARG APP_GID=1001

RUN apt-get update \
  && apt-get install -y --no-install-recommends curl \
  && rm -rf /var/lib/apt/lists/* \
  && groupadd -g "${APP_GID}" appuser \
  && useradd --create-home --no-log-init -u "${APP_UID}" -g "${APP_GID}" appuser

COPY --chown=appuser:appuser --from=builder /app/node_modules ./node_modules
COPY --chown=appuser:appuser --from=builder /app/server.js ./

USER appuser

EXPOSE 3000

ENTRYPOINT ["node", "server.js"]
