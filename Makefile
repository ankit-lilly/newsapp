APP_NAME = newsmaxxing
BUILD_FLAGS = -ldflags="-s -w"


.PHONY: fmt
fmt:
	@echo "Formatting files"
	@go tool templ fmt .
	@gofmt -s -w .

.PHONY: build
build: clean fmt
	@echo "Building $(APP_NAME)..."
	@go tool templ generate
	@tailwindcss -i static/css/style.css -o static/dist/css/style.css --minify
	@bun build static/js/main.js --outdir ./static/dist/js --minify
	@cp -r static/icons ./static/dist/icons 
	@cp static/site.webmanifest ./static/dist/site.webmanifest
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
