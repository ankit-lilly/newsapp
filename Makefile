APP_NAME = newsapp

BUILD_FLAGS = -ldflags="-s -w"

AIR_INSTALL_CMD = go install github.com/cosmtrek/air@latest

.PHONY: install
install:
	@echo "Installing deps"
	@$(AIR_INSTALL_CMD)

.PHONY: run
run:
	@echo "Running ${APP_NAME} in development mode"
	@templ generate
	@air -c .air.toml

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@templ generate
	@go build $(BUILD_FLAGS) -o $(APP_NAME) ./cmd/main.go

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_NAME)

