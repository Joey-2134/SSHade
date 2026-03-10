## Multi-stage build for SSHade
##
## Usage (from repo root):
##   docker build -t sshade .
##   docker run --rm -it -p 2222:2222 -v sshade-data:/data \
##     -e SSHADE_DB_PATH=/data/sshade.db \
##     -e SSHADE_HOST_KEY_PATH=/data/id_ed25519 \
##     sshade

# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /sshade .

# ---- Runtime stage ----
FROM alpine:3.20

WORKDIR /app

RUN adduser -D -H -s /sbin/nologin sshade && \
  apk add --no-cache ca-certificates && \
  mkdir -p /data && \
  chown -R sshade:sshade /data

USER sshade

ENV SSHADE_HOST=0.0.0.0
ENV SSHADE_PORT=2222
ENV SSHADE_DB_PATH=/data/sshade.db
ENV SSHADE_HOST_KEY_PATH=/data/id_ed25519

EXPOSE 2222/tcp

VOLUME ["/data"]

COPY --from=builder /sshade /app/sshade

ENTRYPOINT ["/app/sshade"]

