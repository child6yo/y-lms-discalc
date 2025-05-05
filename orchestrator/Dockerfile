FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

ENV CGO_ENABLED=1
ENV CC=gcc

WORKDIR /orchestrator

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o orchestrator ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /orchestrator

COPY --from=builder /orchestrator/orchestrator .
COPY client/ /orchestrator/client/

RUN mkdir -p /orchestrator/database && \
    chown -R 1000:1000 /orchestrator/database && \
    chmod -R 755 /orchestrator/database

VOLUME /orchestrator/database

CMD ["./orchestrator"]