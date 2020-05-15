build-design:
	$(MAKE) -C site build-design
	cp site/dist/* pkg/webui/content/static/raw

build-content:
	go generate github.com/libesz/poolmanager/pkg/webui/content/static
	go generate github.com/libesz/poolmanager/pkg/webui/content/templates

prepare-site: build-design build-content

build: prepare-site
	go build cmd/poolmanager/main.go

build-pi-zero: prepare-site
	env GOOS=linux GOARCH=arm GOARM=5 go build cmd/poolmanager/main.go

test:
	go test -race -timeout 60s -covermode=atomic -coverprofile=cover.out ./...