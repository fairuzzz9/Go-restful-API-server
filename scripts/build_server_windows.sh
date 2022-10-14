#!/usr/bin/env bash

export GOOS=windows
export GOARCH=amd64
swag init -g cmd/main.go
go build -o ./bin/server_windows ./cmd/