build:
	@echo "Building..."
	@go build -o bin/main cmd/service/main.go

clean:
	@echo "Cleaning"
	@rm -f bin/main

start:
	@bin/main

run:
	@go run cmd/service/main.go

watch:
	@air

inspect:
	@go run cmd/inspect/main.go
