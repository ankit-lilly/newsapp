APP_NAME = newsapp

BUILD_FLAGS = -ldflags="-s -w"

AIR_INSTALL_CMD = go install github.com/air-verse/air@latest

.PHONY: install
install:
	@echo "Installing deps"
	@npm i -D daisyui@latest
	@go install github.com/a-h/templ/cmd/templ@latest
	@$(AIR_INSTALL_CMD)

.PHONY: run
run:
	@echo "Running ${APP_NAME} in development mode"
	@templ generate
	@npx concurrently -k "npx tailwindcss -i assets/css/style.css -o assets/dist/css/style.css --minify --watch" "air -c .air.toml"

.PHONY: fmt
fmt:
	@echo "Formatting files"
	@templ fmt .
	@gofmt -s -w .

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@templ generate
	@npx tailwindcss -i assets/css/style.css -o assets/dist/css/style.css --minify
	@go build $(BUILD_FLAGS) -o $(APP_NAME) ./cmd/main.go

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_NAME)
	@rm -rf assets/dist
	@rm -rf node_modules
	@rm -rf package*
	@rm -rf tmp
