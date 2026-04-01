BACKEND_DIR  := url-shortener-go-backend
FRONTEND_DIR := url-shortener-frontend
BINARY       := server
IMAGE_BACKEND  := url-shortener-backend
IMAGE_FRONTEND := url-shortener-frontend

.PHONY: all build build-backend build-frontend \
        lint lint-backend lint-frontend \
        fmt test clean \
        docker-build docker-build-backend docker-build-frontend

all: lint build

build: build-backend build-frontend

build-backend:
	cd $(BACKEND_DIR) && go build -ldflags="-s -w" -o bin/$(BINARY) ./cmd/server

build-frontend:
	cd $(FRONTEND_DIR) && npm ci && npm run build

lint: lint-backend lint-frontend

lint-backend:
	cd $(BACKEND_DIR) && golangci-lint run ./...

lint-frontend:
	cd $(FRONTEND_DIR) && npm run lint

fmt:
	cd $(BACKEND_DIR) && gofmt -w ./..
	cd $(BACKEND_DIR) && goimports -w ./..

test:
	cd $(BACKEND_DIR) && go test -race -count=1 ./...

clean:
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules

docker-build: docker-build-backend docker-build-frontend

docker-build-backend:
	docker build -t $(IMAGE_BACKEND) ./$(BACKEND_DIR)

docker-build-frontend:
	docker build -t $(IMAGE_FRONTEND) ./$(FRONTEND_DIR)
