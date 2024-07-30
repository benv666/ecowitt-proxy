all: build-local

build-local:
	#go mod download
	go build -o ecoproxy ./cmd/main.go

run: build-local
	./ecoproxy

build-docker:
	go fmt ./...
	docker build -t ecoproxy:latest .

test: build-local
	@./ecoproxy

