VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS  = -ldflags "-X main.version=$(VERSION)"
BINARY   = taq

.PHONY: build install clean

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

clean:
	rm -f $(BINARY)
