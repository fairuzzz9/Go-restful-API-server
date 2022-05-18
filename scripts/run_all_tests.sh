#!/usr/bin/env bash

go clean -testcache
go test -v -cover -race ./...
