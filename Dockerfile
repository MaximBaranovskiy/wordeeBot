FROM golang:1.22.1-bookworm

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o wordee ./cmd/main.go

CMD ["./wordee"]