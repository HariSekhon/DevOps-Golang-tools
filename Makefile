#
#  Author: Hari Sekhon
#  Date: 2013-02-03 10:25:36 +0000 (Sun, 03 Feb 2013)
#
#  https://github.com/harisekhon/devops-golang-tools
#
#  License: see accompanying LICENSE file
#
#  https://www.linkedin.com/in/harisekhon
#

# ===================
# bootstrap commands:

# setup/bootstrap.sh
#
# OR
#
# Alpine:
#
#   apk add --no-cache git make && git clone https://github.com/harisekhon/devops-golang-tools go-tools && cd go-tools && make
#
# Debian / Ubuntu:
#
#   apt-get update && apt-get install -y make git && git clone https://github.com/harisekhon/devops-golang-tools go-tools && cd go-tools && make
#
# RHEL / CentOS:
#
#   yum install -y make git && git clone https://github.com/harisekhon/devops-golang-tools go-tools && cd go-tools && make

# ===================

# would fail bootstrapping on Alpine
#SHELL := /usr/bin/env bash

ifneq ("$(wildcard bash-tools/Makefile.in)", "")
	include bash-tools/Makefile.in
endif

REPO := HariSekhon/DevOps-Golang-tools

# must not be a directory with go.mod
export GOPATH := $(HOME)/go
# use default or allow to be overridden by cross-compiling targets (golang-linux, golang-darwin)
ifndef GOBIN
	export GOBIN  := $(PWD)/bin
endif

CODE_FILES := $(shell find . -type f -name '*.go' | grep -v -e bash-tools -e /lib/ -e /src/)

.PHONY: build
build: init
	@echo =========================
	@echo DevOps Golang Tools Build
	@echo =========================
	@$(MAKE) git-summary

	@#$(MAKE) system-packages-golang
	# provides OS specific build packages for compiling golang modules, such as gcc, and distro specific stdlib packages, build-essential on debian/ubuntu etc
	$(MAKE) system-packages
	@bash-tools/install_packages_if_absent.sh golang

	$(MAKE) golang

.PHONY: init
init:
	git submodule update --init --recursive

.PHONY: golang
golang: golang-version
	./compile.sh
	@echo 'BUILD SUCCESSFUL (go-tools)'

.PHONY: go
go: golang
	@:

.PHONY: golang-mac
golang-mac: golang-darwin
	@:

.PHONY: golang-darwin
golang-darwin:
	GOOS=darwin GOARCH=amd64 GOBIN="$$GOPATH/bin.darwin.amd64" $(MAKE) golang

.PHONY: darwin
darwin: golang-darwin
	@:

# doesn't work yet, issues with -race and also runtime/cgo(__TEXT/__text): relocation target x_cgo_inittls not defined
.PHONY: golang-linux
golang-linux:
	GOOS=linux GOARCH=amd64 GOBIN="$$GOPATH/bin.linux.amd64" CGO_ENABLED=1 $(MAKE) golang

.PHONY: linux
linux: golang-linux
	@:

.PHONY: test-lib
test-lib:
	cd lib && $(MAKE) test

.PHONY: test
test: # test-lib
	tests/all.sh

.PHONY: basic-test
basic-test: test-lib
	bash-tools/check_all.sh

.PHONY: install
install: build
	@echo "No installation needed, just add '$(PWD)' to your \$$PATH"

.PHONY: clean
clean: go-clean
	@rm -vfr bin bin.darwin.amd64 bin.linux.amd64

.PHONY: deep-clean
deep-clean: clean
	@#cd go-lib && $(MAKE) deep-clean
	@echo "Deep cleaning, removing pkg/*"
	@rm -fr pkg/*

# ============================================================================ #
#                       I n d i v i d u a l   B u i l d s
# ============================================================================ #

# Classic Make - individual binaries targets with no-rebuild without clean

# Magic to create dynamic targets
%:: %.go
	@$(MAKE) bin/$@
%.go: %
	@#$(MAKE) $<
# :: means terminal - don't call implicit %.go target, just require the file exists in order for this dynamic target
bin/%:: %.go
	go install -race $<


.PHONY: all
all: colors httpfirst uniq2 welcome
	@:
