install-npm-deps:
	$(MAKE) -C site install-npm-deps

build-site:
	$(MAKE) -C site build-site

build:
	go build cmd/poolmanager/main.go 