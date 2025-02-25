FROM golang:1.23

WORKDIR /agent

COPY go.mod ./

COPY ./ ./

RUN go build -o agent ./cmd/main.go

CMD ["./agent"]