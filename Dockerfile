# Base
FROM golang:1.20.1-alpine AS builder

RUN apk add --no-cache git build-base gcc musl-dev
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build ./cmd/packetengine

FROM alpine:3.17.2
RUN apk -U upgrade --no-cache \
    && apk add --no-cache bind-tools ca-certificates \
    && addgroup packetengine \
    && adduser -D -G packetengine -s /bin/sh packetengine \
    && mkdir -p /home/packetengine/.config/packetengine \
    && chown -R packetengine. /home/packetengine

COPY --from=builder /app/packetengine /usr/local/bin/

WORKDIR /home/packetengine
USER packetengine

VOLUME [ "/home/packetengine/.config/packetengine" ]

ENTRYPOINT ["packetengine"]