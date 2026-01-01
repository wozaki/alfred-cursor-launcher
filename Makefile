.PHONY: build build-universal clean workflow test lint

BINARY_NAME=alfred-cursor-launcher
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=bin
DIST_DIR=dist

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w -X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/alfred-cursor-launcher

build-universal:
	@echo "Building Universal Binary for macOS..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building amd64..."
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-amd64 ./cmd/alfred-cursor-launcher
	@echo "Building arm64..."
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-arm64 ./cmd/alfred-cursor-launcher
	@echo "Creating Universal Binary..."
	lipo -create -output $(BUILD_DIR)/$(BINARY_NAME) \
		$(BUILD_DIR)/$(BINARY_NAME)-amd64 $(BUILD_DIR)/$(BINARY_NAME)-arm64
	@rm $(BUILD_DIR)/$(BINARY_NAME)-amd64 $(BUILD_DIR)/$(BINARY_NAME)-arm64
	@echo "Universal Binary created: $(BUILD_DIR)/$(BINARY_NAME)"

workflow: build-universal
	@echo "Creating Alfred Workflow..."
	@mkdir -p $(DIST_DIR)/workflow
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(DIST_DIR)/workflow/
	@cp workflow/info.plist $(DIST_DIR)/workflow/
	@if [ -f workflow/icon.png ]; then cp workflow/icon.png $(DIST_DIR)/workflow/; fi
	@cd $(DIST_DIR)/workflow && zip -r ../$(BINARY_NAME).alfredworkflow . > /dev/null
	@echo "Alfred Workflow created: $(DIST_DIR)/$(BINARY_NAME).alfredworkflow"

test:
	go test -v ./...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)

