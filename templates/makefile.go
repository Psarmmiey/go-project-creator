package templates

import (
	"fmt"
	"os"
)

var makeFile = `
UNAME_S := $(shell uname -s)

BREW := $(shell command -v brew 2> /dev/null)
GOLANG := $(shell command -v go 2> /dev/null)
JQ := $(shell command -v jq 2> /dev/null)
DOCKER := $(shell command -v docker 2> /dev/null)
STATICCHECK := $(shell command -v staticcheck 2> /dev/null)
MAILHOG := $(shell command -v mailhog 2> /dev/null)

ifeq ($(UNAME_S),Darwin)
## mac
setup:
	$(info "Installing dev dependencies on MacOS")

ifndef BREW
	$(warning "Installing brew")
	/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
endif
	brew --version

ifndef JQ
	$(warning "Installing jq")
	brew install jq
endif
	jq --version

ifndef GOLANG
	$(warning "Installing golang")
	brew install go
endif
	go version

ifndef DOCKER
	$(warning "Installing docker")
	brew install --cask docker
	open /Applications/Docker.app
endif
	docker --version 
ifeq (,$(wildcard .env))
	$(warning "Creating .env")
	cp .env.sample .env
endif

ifndef STATICCHECK
	$(warning "Installing staticcheck")
	brew install staticcheck
endif

ifndef MAILHOG
	$(warning "Installing mailhog")
	brew install mailhog
endif


else
## non mac
	$(warning "Unsupported platform $(UNAME_S)")
endif


build:
	go mod tidy
	staticcheck ./...
	go build 

run_deps:
	/usr/local/opt/mailhog/bin/MailHog

run:
	staticcheck .
	go run main.go

test: 
	staticcheck .
	go test -v ./...

extract-i18n:
	$$HOME/go/bin/goi18n extract
	$$HOME/go/bin/goi18n merge active.en.toml active.fr.toml
	
check-doc:
	go run github.com/Psarmmiey/check-comment@latest
	

`

func CreateMakeFile() {
	// create makefile
	err := os.WriteFile("Makefile", []byte(makeFile), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
