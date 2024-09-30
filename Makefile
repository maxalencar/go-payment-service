.PHONY: run
run:
	@echo "Running..."
	@go run cmd/app/main.go

.PHONY: test
test:
	@echo "Tunning Tests..."
	@go test -v ./...