APP ?= loan
SERVICE_PATH = $(strip ./apps/$(APP))
OUTPUT_BIN = $(SERVICE_PATH)/$(APP)

run:
	@echo "Running $(APP)..."
	@cd $(SERVICE_PATH) && go run main.go

dev:
	@echo "Running $(APP) in development mode..."
	@cd $(SERVICE_PATH) && CompileDaemon -build="go build -o $(APP)-dev" -command="./$(APP)-dev"

build-linux:
	@echo "Building $(APP) for Linux..."
	@cd $(SERVICE_PATH) && GOOS=linux GOARCH=amd64 go build -o $(APP)-linux main.go

build-windows:
	@echo "Building $(APP) for Windows..."
	@cd $(SERVICE_PATH) && GOOS=windows GOARCH=amd64 go build -o $(APP)-windows.exe main.go

build-macos:
	@echo "Building $(APP) for macOS..."
	@cd $(SERVICE_PATH) && GOOS=darwin GOARCH=amd64 go build -o $(APP)-macos main.go

run-linux:
	@echo "Running $(APP) Linux binary..."
	@cd $(SERVICE_PATH) && ./$(APP)-linux

run-windows:
	@echo "Running $(APP) Windows binary..."
	@cd $(SERVICE_PATH) && ./$(APP)-windows.exe

run-macos:
	@echo "Running $(APP) macOS binary..."
	@cd $(SERVICE_PATH) && ./$(APP)-macos

clean:
	@echo "Cleaning up binaries in $(APP)..."
	@rm -f $(SERVICE_PATH)/$(APP)-linux
	@rm -f $(SERVICE_PATH)/$(APP)-windows.exe
	@rm -f $(SERVICE_PATH)/$(APP)-macos
	@rm -f $(SERVICE_PATH)/$(APP)-dev