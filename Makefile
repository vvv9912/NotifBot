BINARY_NAME = NotifBotr
BUILD_DIR = NotifBotrBuild

build:
	@echo "  >  Building binary..."
	@mkdir -p ${BUILD_DIR}
	@GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME} cmd/main.go
	@echo "  >  Moving 'config' folder to build directory..."
	@cp config.hcl  ${BUILD_DIR}/config.hcl

run:
	@echo "  >  Run..."
	@.${BUILD_DIR}/${BINARY_NAME}

clean:
	@echo "  >  Cleaning build cache"
	@go clean
	@rm -rf ${BUILD_DIR}/${BINARY_NAME}