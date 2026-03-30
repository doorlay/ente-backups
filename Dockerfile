FROM golang:1.23-bookworm AS builder
WORKDIR /src
COPY go.mod main.go ./
RUN go build -o /ente-sync .

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
      ca-certificates \
      curl \
      restic && \
    rm -rf /var/lib/apt/lists/*

# Install supercronic
ARG TARGETARCH
RUN curl -fsSL "https://github.com/aptible/supercronic/releases/download/v0.2.33/supercronic-linux-${TARGETARCH}" \
      -o /usr/local/bin/supercronic && \
    chmod +x /usr/local/bin/supercronic

# Install ente CLI
RUN curl -fsSL "https://github.com/ente-io/ente/releases/latest/download/ente-linux-${TARGETARCH}" \
      -o /usr/local/bin/ente && \
    chmod +x /usr/local/bin/ente

COPY --from=builder /ente-sync /usr/local/bin/ente-sync
COPY restic-backup.sh /usr/local/bin/restic-backup.sh
COPY crontab /etc/crontab

RUN mkdir -p /srv/backups/tmp /data/backups

ENV TMPDIR=/srv/backups/tmp
ENV TZ=America/Los_Angeles

CMD ["supercronic", "/etc/crontab"]
