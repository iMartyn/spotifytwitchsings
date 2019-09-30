# DISCLAIMER: this Makefile has only been tested on Linux with git, head and cut installed.
# It _should_ work on other OS' but the intention is that it should be run on CI.

PKGS := $(shell go list ./... | grep -v /vendor)
BINARY := codetest

GOPATH ?= $(HOME)/go
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

.PHONY: all
all: deps test $(BINARY)

.PHONY: deps
deps:
	echo '$(GOPATH)'
	go get github.com/spf13/cobra/cobra
	go get github.com/frankban/quicktest
	go get github.com/gosuri/uiprogress
	go get github.com/chrisvdg/spotify
	sh -c 'cd $(GOPATH)/src/github.com/chrisvdg/spotify && git checkout all_pages && git pull'
	go get github.com/gorilla/mux
	go get google.golang.org/api/googleapi/transport

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter . src


.PHONY: test
test: lint
	go test $(PKGS)

.PHONY: clean
clean:
	rm -f ${BINARY}
	rm -rf release
	go clean

.PHONY: build
$(BINARY): deps
	go build -o ${BINARY}

PLATFORMS := linux darwin windows
os = $(word 1, $@)
GITHASH := $(shell git log --oneline | head -n1 | cut -f1 -d' ')
VERSION ?= git-${GITHASH}

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	echo $(os) ${VERSION} ${BINARY}
	# We need to run this here because different OS' have different deps
	GOOS=$(os) make deps
	GOARCH=amd64 GOOS=$(os) go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: release
release: $(PLATFORMS)

.PHONY: docker
docker:
	docker build .