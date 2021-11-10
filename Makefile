#!/bin/bash

BUILD_TARGETS=windows linux

APPID=edgex-datamonitor

FYNECROSS=fyne-cross

# BIN_NAME=edgex-datamonitor


#	../../$(FYNECROSS) $(1) -app-id $(APPID) -arch=amd64  -env="CGO_CFLAGS=-I/usr/include -I/usr/include/x86_64-linux-gnu" -env="CGO_LDFLAGS=-L/usr/lib/x86_64-linux-gnu" -debug=true ./cmd/app

define compile_target
	cd ./src/edgex-foundry-datamonitor && \
	$(FYNECROSS) $(1) -app-id $(APPID) -arch=amd64 ./cmd/app


endef

PHONY: install-deps
install-deps:
	go install github.com/fyne-io/fyne-cross@latest




build-builder:
	cd ./src/fyne-cross-zq && make fyne-cross-zq

darwin-compile:
	@echo "\n❗❗❗❗ 🍎"
	@echo "OSX/Darwin/Apple cross-compiling requires a darwin host and/or some manual steps along with the acceptance of Xcode license terms\n"
	@echo "Please follow the link below:"
	@echo https://github.com/fyne-io/fyne-cross#build-the-docker-image-for-osxdarwinapple-cross-compiling
	@echo "\n❗❗❗❗\n"


cross-compile: $(BUILD_TARGETS)
	make darwin-compile

refresh-windows:
	cd ./src/edgex-foundry-datamonitor && \
	$(FYNECROSS) windows -app-id $(APPID) -arch=amd64 ./cmd/app
	make run-windows

refresh-linux:
	cd ./src/edgex-foundry-datamonitor && \
	$(FYNECROSS) linux -app-id $(APPID) -arch=amd64 ./cmd/app
	make run-linux

run-windows:
	./src/edgex-foundry-datamonitor/fyne-cross/bin/windows-amd64/edgex-foundry-datamonitor.exe

run-linux:
	./src/edgex-foundry-datamonitor/fyne-cross/bin/linux-amd64/edgex-foundry-datamonitor

$(BUILD_TARGETS):
	$(call compile_target,$(@))