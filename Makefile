BINARY_NAME=voiceops-pi

ifneq (,$(wildcard .env))
	include .env
	export
endif

run:
	go run cmd/voiceops/main.go

build-pi:
	@echo "Building for Raspberry Pi..."
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME) ./cmd/voiceops/main.go

deploy-pi: build-pi
	@if [ -z "$(PI_HOST)" ] || [ -z "$(PI_USER)" ] || [ -z "$(PI_KEY)" ]; then \
		echo "Error: PI_HOST or PI_USER is not set. Check your .env file"; \
		exit 1; \
	fi

	@echo "Deploying to Raspberry Pi at $(PI_USER)@$(PI_HOST)..."

	scp -i $(PI_KEY) -o StrictHostKeyChecking=no $(BINARY_NAME) $(PI_USER)@$(PI_HOST):$(PI_DIR)/$(BINARY_NAME)_new
	
	ssh -i $(PI_KEY) -o StrictHostKeyChecking=no $(PI_USER)@$(PI_HOST) "mv $(PI_DIR)/$(BINARY_NAME)_new $(PI_DIR)/$(BINARY_NAME) && chmod +x $(PI_DIR)/$(BINARY_NAME) && sudo systemctl restart voiceops"
	
	@echo "Deployed and restarted!"

clean:
	go clean
	rm $(BINARY_NAME)