.PHONY: build
build:
	go build -o app ./cmd/main.go

.PHONY: run
run:
	go run ./cmd

.PHONY: test
test:
	go test ./... -v

.PHONY: tidy
tidy:
	go get -d ./...
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run -c configs/linter/.golangci.yml -v --fix

.PHONY: build-docker-image
build-docker-image:
	docker build -t registry.gitlab.com/mfcekirdek/kv-store:latest -f deployments/Dockerfile .