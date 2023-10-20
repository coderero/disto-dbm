NAME="Auhentication"
VERSION="0.0.1"

build:
	@echo "\033[92mBuilding $(NAME) $(VERSION)\033[0m"
	@go build -o ./bin/server -ldflags "-X main.version=$(VERSION)" ./cmd/server/main.go

run: build
	@echo "\033[92mRunning $(NAME) $(VERSION)\033[0m"
	@./bin/server