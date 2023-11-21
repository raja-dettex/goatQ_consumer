run: build
	@./bin/goatQ_consumer
build:
	@go build -o ./bin/goatQ_consumer