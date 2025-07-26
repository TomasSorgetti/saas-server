FROM golang:1.24.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate

COPY . ./

RUN go build -o main ./cmd/api

CMD ["./main"]