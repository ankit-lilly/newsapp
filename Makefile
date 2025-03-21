APP_NAME = newsapp
BUILD_FLAGS = -ldflags="-s -w"
AIR_INSTALL_CMD = go install github.com/air-verse/air@latest

.PHONY: install
install:
	@echo "Installing deps"
	@go install github.com/a-h/templ/cmd/templ@latest
	@bun install
	@$(AIR_INSTALL_CMD)

.PHONY: run
run:
	@echo "Running ${APP_NAME} in development mode"
	@templ generate
	@bun build static/js/*.js --outdir ./static/dist/js --minify
	@cp static/icons ./static/dist/icons -r
	@cp static/site.webmanifest ./static/dist/site.webmanifest
	@bunx concurrently -k "bunx @tailwindcss/cli -i static/css/style.css -o static/dist/css/style.css --minify --watch" "air -c .air.toml"

.PHONY: fmt
fmt:
	@echo "Formatting files"
	@templ fmt .
	@gofmt -s -w .

generate:
	@echo "Generating files"
	@templ generate
	@bunx @tailwindcss/cli -i static/css/style.css -o static/dist/css/style.css --minify
	@bun build static/js/*.js --outdir ./static/dist/js --minify

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@templ generate
	@bunx @tailwindcss/cli -i static/css/style.css -o static/dist/css/style.css --minify
	@bun build static/js/*.js --outdir ./static/dist/js --minify
	@cp static/icons ./static/dist/icons -r
	@cp static/site.webmanifest ./static/dist/site.webmanifest
	@go build $(BUILD_FLAGS) -o $(APP_NAME) main.go

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_NAME)
	@rm -rf static/dist
	@rm -rf node_modules
	@rm -rf tmp

