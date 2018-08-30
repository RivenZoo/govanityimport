# usage:
# 1. put make_gopath.mk to your project top directory	
# 2. add `include make_gopath.mk` to your project Makefile
# 3. make sure target `gopath` to be called
# 4. `cd $(PROJ_GOPATH)` then build
# 5. set `REPO_ROOT` to specify repo dir in your go import path, e.g. github.com/yourusername

.PHONY: gopath

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))

OS = $(shell uname)
GOENV_HOME = $(HOME)/.goenv
PROJ_DIR = $(dir $(mkfile_path))
PROJ_NAME = $(shell basename $(PROJ_DIR))
REPO_ROOT = github.com/RivenZoo

ifneq ($(REPO_ROOT),)
PROJ_NAME := $(REPO_ROOT)/$(PROJ_NAME)
endif

PROJ_GOPATH =
GOPATH =

ifeq ($(OS),Linux)
pathmd5 = $(shell echo $(PROJ_NAME) | md5sum | awk '{print $1}')
endif
ifeq ($(OS),Darwin)
pathmd5 = $(shell echo $(PROJ_NAME) | md5)
endif
ifeq ($(pathmd5),)
$(error GOPATH not generate on $(OS))
endif

GOPATH = ${GOENV_HOME}/$(shell echo $(pathmd5) | cut -c1-16)
PROJ_GOPATH = $(GOPATH)/src/$(PROJ_NAME)

$(GOPATH): | $(GOPATH)/workdir

$(GOPATH)/workdir:
	mkdir -p $(GOPATH) $(GOPATH)/bin $(GOPATH)/src $(GOPATH)/pkg
	mkdir -p $(dir $(PROJ_GOPATH))
	echo $(PROJ_DIR) > $(GOPATH)/workdir
	ln -s $(PROJ_DIR) $(PROJ_GOPATH)

gopath: $(GOPATH)


export GOPATH
export PROJ_GOPATH
