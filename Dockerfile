# syntax=docker/dockerfile:1
FROM golang:1.22.5
WORKDIR /vezbamo
COPY . .
RUN go mod download
RUN go build -o bin src/main.go
EXPOSE 10000
ENTRYPOINT [ "bin/main" ]
