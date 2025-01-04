#!/bin/bash

export PKG_CONFIG_PATH=$PWD

go env -w GOPROXY=https://goproxy.cn,direct

gofmt -s -w .
golangci-lint run

go mod tidy
go mod verify
