FROM golang:1.22.1-bookworm

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o wordee ./cmd/main.go

ENV DATABASE_URL="host=db port=5432 user=USER password=PASSWORD dbname=DBNAME sslmode=disable"

CMD ["./wordee","&&","goose","-dir","./internal/migrations","postgres","${DATABASE_URL}","up"]