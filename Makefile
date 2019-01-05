.PHONY: dep
dep:
	@ go mod download

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

.PHONY: cover
cover:
	@ go test -coverprofile cover.out ./...
	@ go tool cover -html=cover.out -o cover.html
	@ rm -f cover.out

.PHONY: build
build:
	@ go build -o ./bin/nott ./cmd/nott

.PHONY: run
run:
	@ ./bin/nott

.PHONY: docker
docker:
	@ docker image build -t tetafro/nott-backend-go .
