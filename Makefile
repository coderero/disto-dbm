NAME="Auhentication"
VERSION="0.0.1"

build:
	@echo "\033[92mBuilding $(NAME) $(VERSION)\033[0m"
	@go build -o ./bin/server -tags=jsoniter,netgo  -ldflags '-s -w' ./cmd/main.go

run: 
	@echo "\033[92mRunning $(NAME) $(VERSION)\033[0m"
	@go run -tags=jsoniter,netgo ./cmd/main.go