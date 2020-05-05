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

CODE_FILES := $(shell find . -type f -name '*.go' | grep -v -e bash-tools -e /lib/)

.PHONY: build
build: init golang-version
	@echo =========================
	@echo DevOps Golang Tools Build
	@echo =========================
	@$(MAKE) git-summary

	if [ -z "$(CPANM)" ]; then make; exit $$?; fi
	@#$(MAKE) system-packages-golang

	$(MAKE) golang

.PHONY: init
init:
	git submodule update --init --recursive

.PHONY: golang
golang: golang-version
	@for x in *.go; do \
		echo "go build $$x"; \
		go build "$$x"; \
		echo; \
	done
	@echo 'BUILD SUCCESSFUL (go-tools)'

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

#.PHONY: clean
#clean:
#	cd lib && $(MAKE) clean

#.PHONY: deep-clean
#deep-clean: clean
#	cd go-lib && $(MAKE) deep-clean
