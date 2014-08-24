
GOPATH = /tmp/go.build.dsapid

SERVER_BIN = dsapid
REPOSITORY_BASE = github.com/MerlinDMC

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
