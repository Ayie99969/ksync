GO_VERSION := $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1,2)

###############################################################################
###                                  Build                                  ###
###############################################################################

build: ensure_version
	go build -mod=readonly -o ./build/ksync ./cmd/ksync/main.go

###############################################################################
###                                 Checks                                  ###
###############################################################################

ensure_version:
ifneq ($(GO_VERSION),1.21)
	$(error ❌  Please run Go v1.21.x..)
endif