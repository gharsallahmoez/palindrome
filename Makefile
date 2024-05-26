VERSION=1.0.0

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
GORUN=$(GOCMD) run
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod
SERVER_PATH=./cmd
DOCKER_IMAGE_TAG=gharsallahmoez/messages
BINARY_SERVER_NAME=messagessvc
DOCKER_PORT= 8080

init:
	$(GOMOD) download

run:
	cd ${SERVER_PATH} && $(GORUN) .

build:
	$(GOBUILD) -o ./bin/messages/$(BINARY_SERVER_NAME) -v $(SERVER_PATH)

tidy:
	$(GOMOD) tidy

tool-moq:
	$(GOINSTALL) github.com/matryer/moq@v0.3.4

moq:
	moq -out server/http/zmoq_infra_database_test.go -pkg http_test infra/database Database

test:
	$(GOTEST) -v ./...

test-race:
	$(GOTEST) --race -v ./...

test-e2e:
	$(GOTEST) -v e2e/e2e_test.go

test-coverage:
	$(GOTEST)  ./... -coverprofile cover.out.tmp
	cat cover.out.tmp | grep -v "main.go" > cover.out
	go tool cover -func cover.out

tool-lint:
	$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1

lint:
	golangci-lint run

govulncheck: ## Check dependencies for vulnerabilities.
	@command -v govulncheck > /dev/null || $(GOINSTALL) golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

clean-test-cache:
	$(GOCLEAN) -testcache

clean-mod-cache:
	$(GOCLEAN) -modcache

docker-build:
	docker build -t ${DOCKER_IMAGE_TAG} .

docker-run:
	docker run -p 8080:8080 ${DOCKER_IMAGE_TAG}

docker-compose-run:
	docker-compose up -d --build

docker-compose-down:
	docker-compose down
