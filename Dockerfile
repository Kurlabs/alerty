FROM golang:1.13-alpine as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/websites-cron cmd/websites-cron/run.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/sockets-cron cmd/sockets-cron/run.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/robots-cron cmd/robots-cron/run.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/brain cmd/brain/brain.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/controller cmd/controller/controller.go

# Final stage
FROM scratch

COPY --from=builder /app/bin /usr/bin

WORKDIR /app
