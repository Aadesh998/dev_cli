APP_NAME=dev_cli
BUILD_DIR := ./bin
ENTRY=main.go
CONFIG_FILE=/home/aadesh-kumar/.config/dev_cli

WIN_EXE := $(BUILD_DIR)/$(APP_NAME).exe
LINUX_BIN := $(BUILD_DIR)/$(APP_NAME)

.PHONY: all build run build_win build_linux clean

all: build

build: build_linux build_win 

build_win:
	@echo "Building Windows executable $(WIN_EXE)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(WIN_EXE) $(ENTRY)

build_linux:
	@echo "Building Linux executable $(LINUX_BIN)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_BIN) $(ENTRY)

run:
	go run $(ENTRY)

remove_config:
	@echo "Removing config file"
	@rm -rf $(CONFIG_FILE)

clean:
	@echo "Cleaning up"
	@rm -rf $(BUILD_DIR)/*