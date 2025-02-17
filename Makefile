TEST_PATH=./...

VERSION := v1.0.0
DATE := $$(date)

MOCK_SERVER_SRC=./internal/server/app/app.go
MOCK_SERVER_DST=./internal/server/app/mocks/mocks.go

MOCK_SERVICE_SRC=./internal/server/service/service.go
MOCK_SERVICE_DST=./internal/server/service/mocks/mocks.go

MOCK_STORAGE_SRC=./internal/server/storage/db/db.go
MOCK_STORAGE_DST=./internal/server/storage/db/mocks/mocks.go

MOCK_GRPC_SRC=./proto/keeper_grpc.pb.go
MOCK_GRPC_DST=./proto/mocks/mocks.go

MOCK_LOGGER_SRC=./internal/server/logger/logger.go
MOCK_LOGGER_DST=./internal/server/logger/mocks/mocks.go

.PHONY=mock-gen

mock-gen:
	$(GOPATH)/bin/mockgen -source=$(MOCK_STORAGE_SRC) -destination=$(MOCK_STORAGE_DST)
	$(GOPATH)/bin/mockgen -source=$(MOCK_GRPC_SRC) -destination=$(MOCK_GRPC_DST)
	$(GOPATH)/bin/mockgen -source=$(MOCK_LOGGER_SRC) -destination=$(MOCK_LOGGER_DST)
	$(GOPATH)/bin/mockgen -source=$(MOCK_SERVER_SRC) -destination=$(MOCK_SERVER_DST)
	$(GOPATH)/bin/mockgen -source=$(MOCK_SERVICE_SRC) -destination=$(MOCK_SERVICE_DST)


.PHONY=lint
lint:
	$(GOPATH)/bin/golangci-lint run

.PHONY=docker-run
docker-run:
	docker compose up --build --abort-on-container-exit

GO_IMPORTS = $(GOPATH)/bin/goimports
GO_FILES = $(shell find . -name '*.go')

.PHONY=make-import
make-import:
	$(GO_IMPORTS) -local github.com/Sofja96/GophKeeper.git -w $(GO_FILES)

.PHONY=build
build:
	go build -o client -ldflags="-X 'github.com/Sofja96/GophKeeper.git/shared/buildinfo.Version=${VERSION}' -X 'github.com/Sofja96/GophKeeper.git/shared/buildinfo.BuildDate=${DATE}'" cmd/client/main.go
	go build -o server -ldflags="-X 'github.com/Sofja96/GophKeeper.git/shared/buildinfo.Version=${VERSION}' -X 'github.com/Sofja96/GophKeeper.git/shared/buildinfo.BuildDate=${DATE}'" cmd/server/main.go


.PHONY: test
test:
	go test -v $(TEST_PATH)

.PHONY: test-without-pb
test-without-pb:
	go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./...
	grep -v ".pb.go" coverage.out > coverage_filtered.out
	go tool cover -func=coverage_filtered.out

.PHONY: test-with-coverage
test-with-coverage:
	go test -coverprofile=cover.out -v $(TEST_PATH)
	make --silent coverage

.PHONY: coverage
coverage:
	go tool cover -html cover.out -o cover.html
	open cover.html

.PHONY: total-coverage
total-coverage:
	go tool cover -func cover.out

.PHONY:clean
clean:
	rm -rf cover.out.html cover.out
