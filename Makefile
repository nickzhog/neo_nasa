test:
	@go test ./...
clean:
	@find . -type f -name "*.log" -delete