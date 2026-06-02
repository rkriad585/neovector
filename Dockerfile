# ── Build stage ────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION
ARG COMMIT
ARG PUBLISHER_NAME=rkriad585
ARG PUBLISHER_EMAIL=rkriad585@gmail.com

RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-s -w \
    -X main.Version=${VERSION} \
    -X main.Commit=${COMMIT} \
    -X main.PublisherName=${PUBLISHER_NAME} \
    -X main.PublisherEmail=${PUBLISHER_EMAIL}" \
    -o /usr/local/bin/neovector .

# ── Runtime stage ──────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /usr/local/bin/neovector /usr/local/bin/neovector

ENV USER=neovector
RUN adduser -D -h /home/neovector neovector && \
    mkdir -p /home/neovector/.config/neostore/neovector && \
    chown -R neovector:neovector /home/neovector

USER neovector
WORKDIR /home/neovector

ENTRYPOINT ["neovector"]
CMD ["--help"]
