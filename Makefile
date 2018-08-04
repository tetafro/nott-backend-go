.PHONY: init
init:
	@ dep ensure -v -vendor-only

.PHONY: dep
dep:
	@ dep ensure -v

.PHONY: mock
mock:
	@ mockgen \
		-source=internal/auth/tokens_repo.go \
		-destination=internal/auth/tokens_repo_mock.go \
		-package=auth
	@ mockgen \
		-source=internal/auth/users_repo.go \
		-destination=internal/auth/users_repo_mock.go \
		-package=auth
	@ mockgen \
		-source=internal/folders/repo.go \
		-destination=internal/folders/repo_mock.go \
		-package=folders
	@ mockgen \
		-source=internal/notepads/repo.go \
		-destination=internal/notepads/repo_mock.go \
		-package=notepads
	@ mockgen \
		-source=internal/notes/repo.go \
		-destination=internal/notes/repo_mock.go \
		-package=notes

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
