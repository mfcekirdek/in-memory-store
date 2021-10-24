.PHONY: build
build:
	go build -o app ./cmd/main.go

.PHONY: run
run:
	go run ./cmd

.PHONY: test
test:
	sh mockgen.sh
	go test ./... -v -coverprofile=unit_coverage.out -tags=unit

.PHONY: unit-coverage-html
unit-coverage-html:
	make test
	go tool cover -html=unit_coverage.out -o unit_coverage.html

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

.PHONY: generate-mocks
generate-mocks:
	sh mockgen.sh