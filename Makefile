# Defining necessary environment variables
export PROJECT_ROOT := $(shell pwd)
export SERVER_WINDOW_SIZE_IN_SECONDS := 60
export SERVER_PORT := :8080

.PHONY: start-server clean

start-server:
	@echo "Starting Server..."
	@go run cmd/main.go

clean:
	@rm -f .hits .windowmetadata