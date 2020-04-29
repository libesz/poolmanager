build-design:
	$(MAKE) -C site build-design

build: build-design
	cp site/dist/* pkg/webui/content/static/raw
	go generate github.com/libesz/poolmanager/pkg/webui/content/static
	go generate github.com/libesz/poolmanager/pkg/webui/content/templates
	go build cmd/poolmanager/main.go
