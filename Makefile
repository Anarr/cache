# Makefile for the Go cache library

.PHONY: test

# Run all tests in the cache library
test:
	go test ./... -v
