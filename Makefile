# Makefile for User Management API

GO_VERSION := 1.24
BINARY_NAME := user-management

.PHONY: test-integration
test-integration:
	go test -v -race -coverprofile=coverage-integration.out --coverpkg=./... ./tests/...
	go tool cover -func=coverage-integration.out
