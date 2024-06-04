BINARY_NAME=goravelApp

build:
	@go mod vendor
	@echo "Building my goravel-demo-app..."
	@go build -o out/${BINARY_NAME} .
	@echo "goravel-demo-app built!"

run: build
	@echo "Starting my goravel-demo-app..."
	@./out/${BINARY_NAME} &
	@echo "My goravel-demo-app started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm out/${BINARY_NAME}
	@echo "Cleaned!"

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run

stop:
	@echo "Stopping my goravel-demo-app..."
	@-pkill -SIGTERM -f "./out/${BINARY_NAME}"
	@echo "Stopped my goravel-demo-app!"

restart: stop start