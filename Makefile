
GOPATH = /tmp/go.build.dsapid

SERVER_BIN = dsapid
REPOSITORY_BASE = github.com/MerlinDMC

VERSION ?= $(shell git describe --abbrev=0 --tags --match 'v*' | sed -e 's/^v//')
COMMIT := $(shell git rev-parse --short HEAD)
OWNER ?= merlindmc
GOOS := linux


all:	$(SERVER_BIN)

$(GOPATH):
	mkdir -p $(GOPATH)
	GOPATH=$(GOPATH) go get -d

$(GOPATH)/src/$(REPOSITORY_BASE)/dsapid:	$(GOPATH)
	mkdir -p $(GOPATH)/src/$(REPOSITORY_BASE)
	ln -sf $(shell pwd) $(GOPATH)/src/$(REPOSITORY_BASE)/dsapid

$(SERVER_BIN):	$(GOPATH)/src/$(REPOSITORY_BASE)/dsapid
	GOPATH=$(GOPATH) go get -d $(REPOSITORY_BASE)/dsapid/server
	GOPATH=$(GOPATH) go build -o $(SERVER_BIN) $(REPOSITORY_BASE)/dsapid/server

clean:
	rm -rf $(GOPATH)
	rm -f $(SERVER_BIN)

release_docker:	$(GOPATH)/src/$(REPOSITORY_BASE)/dsapid
	GOPATH=$(GOPATH) GOOS=$(GOOS) go get -d $(REPOSITORY_BASE)/dsapid/server
	GOPATH=$(GOPATH) GOOS=$(GOOS) go build -o dsapid $(REPOSITORY_BASE)/dsapid/server
	docker build --build-arg dsapid_version=$(VERSION) --build-arg dsapid_commit=$(COMMIT) -t $(OWNER)/dsapid:$(VERSION) -f Dockerfile.scratch .
	docker tag $(OWNER)/dsapid:$(VERSION) $(OWNER)/dsapid:latest
	rm -f dsapid
