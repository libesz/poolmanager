build-site:
	$(MAKE) -C site build-site
	cp -r site/dist/* pkg/webui/content/raw

build-site-dev:
	$(MAKE) -C site build-site-dev
	cp -r site/dist/* pkg/webui/content/raw

build-content:
	go generate github.com/libesz/poolmanager/pkg/webui/content/

prepare-site-dev: build-site-dev build-content

prepare-site: build-site build-content

build: prepare-site-dev
	go build cmd/poolmanager/main.go

build-pi-zero: prepare-site
	env GOOS=linux GOARCH=arm GOARM=5 go build cmd/poolmanager/main.go

test:
	go test -race -timeout 60s -covermode=atomic -coverprofile=cover.out ./...

lint:
	golangci-lint run

validate-code: lint test