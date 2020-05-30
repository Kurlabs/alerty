FROM golang:1.13-alpine as builder

ENV GO111MODULE=on

WORKDIR /alerty

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN touch .env

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/websites-cron cmd/websites-cron/run.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/sockets-cron cmd/sockets-cron/run.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/robots-cron cmd/robots-cron/run.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/brain cmd/brain/brain.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/controller cmd/controller/controller.go

# Final stage
FROM golang:1.13-alpine

COPY --from=builder /alerty/bin /usr/bin
COPY --from=builder /alerty/.env /alerty/.env

WORKDIR /alerty
