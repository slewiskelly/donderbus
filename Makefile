GOVERSION ?= $(shell go list -m -f '{{.GoVersion}}')

.PHONY: build
build: clean
	@CGO_ENABLED=0 go build -o ./bin/ ./cmd/*

.PHONY: build/image
build/image:
	@docker buildx build . -t slewiskelly/donderbus --build-arg GOVERSION=$(GOVERSION)

.PHONY: clean
clean:
	@rm -rf ./bin ./dist

.PHONY: test
test:
	@go test -race ./...

.PHONY: vet
vet:
	@go vet ./...
