.PHONY: init
init:
	@ dep ensure -v -vendor-only

.PHONY: dep
dep:
	@ dep ensure -v

.PHONY: mock
mock:
	@ mockgen \
		-source=internal/storage/repositories.go \
		-destination=internal/storage/repositories_mock.go \
		-package=storage

.PHONY: lint
lint:
	@ golangci-lint run

.PHONY: test
test:
	@ go test ./...

.PHONY: build
build:
	@ go build -o ./bin/nott ./cmd/nott

.PHONY: run
run:
	@ ./bin/nott

.PHONY: docker
docker:
	@ docker build -t tetafro/nott-backend-go .
