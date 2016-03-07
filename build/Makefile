#
# Makefile for kii-cli
#
# Requirement:
#   CentOS 6.x 64bit
#   https://github.com/jordansissel/fpm
#
# Steps:
#   make install
#   make clean
#   make build
#
#   sudo rpm -qpil <NAME.rpm>
#   sudo yum localinstall <NAME.rpm>

RPM_NAME := kii-cli
RPM_VERSION := 0.1.9
RPM_URL := https://github.com/tmtk75/kii-cli
RPM_DESCRIPTION := KiiCloud command line interface
RPM_DEST_DIR := ./
BUILDROOT := ./BUILDROOT
ARC := $(shell uname -i)
SHELL := /bin/bash

install:
	sudo yum -y install gcc ruby ruby-devel ruby-libs rubygems rpm-build
	sudo gem install fpm

build:
	mkdir -p $(BUILDROOT)/usr/local/bin
	wget -nv https://github.com/tmtk75/kii-cli/releases/download/v$(RPM_VERSION)/kii-cli_linux_amd64.gz
	gunzip kii-cli_linux_amd64.gz
	mv kii-cli_linux_amd64 $(BUILDROOT)/usr/local/bin/$(RPM_NAME)
	chmod 755 $(BUILDROOT)/usr/local/bin/$(RPM_NAME)
	fpm -s dir -t rpm \
		-n $(RPM_NAME) \
		-v $(RPM_VERSION) \
		-C $(BUILDROOT) \
		-a $(ARC) \
		--url "$(RPM_URL)" \
		--description "$(RPM_DESCRIPTION)" \
		$(RPM_DEST_DIR)
	rpm -qpil $(RPM_DEST_DIR)$(RPM_NAME)-$(RPM_VERSION)-1.$(ARC).rpm

copy2repo:
	scp $(RPM_DEST_DIR)$(RPM_NAME)-$(RPM_VERSION)-1.$(ARC).rpm $(TARGET_HOSTNAME):/tmp/
	@echo -e "please copy from /tmp/$(RPM_NAME)-$(RPM_VERSION)-1.$(ARC).rpm to /mnt/ebs/rpm-repo/external/x86_64/"

clean:
	rm -rf BUILDROOT
	rm -f *.rpm
