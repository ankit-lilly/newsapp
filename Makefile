APP_NAME = newsmaxxing
BUILD_FLAGS = -ldflags="-s -w"
AIR_INSTALL_CMD = go install github.com/air-verse/air@latest
TAILWIND = static/css/tailwindcss

.PHONY: install
install:
	@echo "Installing dependencies"
	@$(AIR_INSTALL_CMD)
	@go run ./cmd/install

.PHONY: run
run: prepare-static
	@if [ ! -x $(TAILWIND) ]; then \
		echo "Run 'make install' first to download dependencies"; \
		exit 1; \
	fi
	@echo "Running $(APP_NAME) in development mode"
	@$(TAILWIND) -i static/css/style.css -o static/dist/css/style.css --minify --watch & \
		tailwind_pid=$$!; \
		trap 'kill $$tailwind_pid 2>/dev/null || true' INT TERM EXIT; \
		air -c .air.toml

.PHONY: fmt
fmt:
	@echo "Formatting files"
	@go tool templ fmt .
	@gofmt -s -w .

.PHONY: generate
generate: templates css static-assets

.PHONY: templates
templates:
	@echo "Generating templates"
	@go tool templ generate

.PHONY: css
css:
	@if [ ! -x $(TAILWIND) ]; then \
		echo "Run 'make install' first to download dependencies"; \
		exit 1; \
	fi
	@mkdir -p static/dist/css
	@$(TAILWIND) -i static/css/style.css -o static/dist/css/style.css --minify

.PHONY: static-assets
static-assets:
	@if [ ! -f static/js/smd.min.js ]; then \
		echo "Run 'make install' first to download dependencies"; \
		exit 1; \
	fi
	@mkdir -p static/dist/js
	@cp static/js/main.js static/js/smd.min.js static/dist/js/
	@mkdir -p static/dist/icons
	@cp -R static/icons/. static/dist/icons/
	@cp static/site.webmanifest static/dist/site.webmanifest

.PHONY: prepare-static
prepare-static: templates static-assets

.PHONY: build
build: generate
	@echo "Building $(APP_NAME)..."
	@go build $(BUILD_FLAGS) -o $(APP_NAME) main.go

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_NAME)
	@rm -rf static/dist
	@rm -rf tmp

.PHONY: clean-all
clean-all: clean
	@rm -f static/css/tailwindcss
	@rm -f static/css/daisyui.mjs static/css/daisyui-theme.mjs
	@rm -f static/js/smd.min.js
