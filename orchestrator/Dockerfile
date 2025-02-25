FROM golang:1.23

WORKDIR /orchestrator

COPY go.mod ./

COPY ./ ./

RUN go build -o orchestrator ./cmd/main.go

CMD ["./orchestrator"]