#!/usr/bin/env bash

export GOOS=linux
export GOARCH=amd64
swag init -g cmd/main.go
go build -o ./bin/server_linux ./cmd/