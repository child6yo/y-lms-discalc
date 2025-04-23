FROM golang:1.23

WORKDIR /agent

COPY go.mod go.sum ./
COPY ./ ./

RUN go mod tidy

RUN go build -o agent ./cmd/main.go

CMD ["./agent"]