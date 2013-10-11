PROJECT_DIR=$(shell pwd)
GOPATH=$(PROJECT_DIR)/libs:$(PROJECT_DIR)

# https://github.com/webrocket/webrocket/blob/master/Makefile

all: build clean

help:
	@echo "Please use 'make <target>' where <target> is one of"
	@echo "  build    to build :)"

clean:
	rm -rf $(PROJECT_DIR)/libs/pkg
	rm -rf $(PROJECT_DIR)/libs/src/*

build:
	GOPATH=$(GOPATH) go build sync2async

fmt:
	GOPATH=$(GOPATH) go fmt sync2async

doc:
	GOPATH=$(GOPATH) go doc sync2async

run:
	GOPATH=$(GOPATH) go run
