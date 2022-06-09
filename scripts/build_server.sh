#!/usr/bin/env bash
swag init -g cmd/main.go
go build -o ./bin/server_mac ./cmd/