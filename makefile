run:
	@echo "Running Candle..."
	@go run main.go

build:
	@echo "Building Candle..."
	go build -o bin/candle main.go
	@echo "Build succeeded. Binary at bin/candle"

clean:
	@echo "Cleaning Candle..."
	rm -rf bin/
	@echo "Clean succeeded."

test:
	@if [ -z "$$PACKAGE" ]; then \
		go test -v ./...; \
	else \
		go test -v ./$$PACKAGE; \
	fi
