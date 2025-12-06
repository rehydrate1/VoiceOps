BINARY_NAME=voiceops-pi

run:
	go run cmd/voiceops/main.go

build-pi:
	@echo "Building for Raspberry Pi..."
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME) ./cmd/voiceops/main.go

clean:
	go clean
	rm $(BINARY_NAME)