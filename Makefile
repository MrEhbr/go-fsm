SHELL := /usr/bin/env bash -o pipefail
GOPKG ?= github.com/MrEhbr/go-fsm/v2
DOCKER_IMAGE ?=	mrehbr/go-fsm
GOBINS ?= cmd/go-fsm

include rules.mk
