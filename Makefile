NAME="Auhentication"
VERSION="0.0.1"

build:
	@echo "\033[92mBuilding $(NAME) $(VERSION)\033[0m"
	@go build -o ./bin/server -tags=jsoniter,netgo  -ldflags '-s -w -X main.version=${VERSION}' ./cmd/main.go

run: build
	@echo "\033[92mRunning $(NAME) $(VERSION)\033[0m"
	@./bin/server