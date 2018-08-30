.PHONY: all main gopath proto web

PROTOC = protoc
PB_DIR = .
PROTO_DIR = api_define
PROJNAME = $(notdir $(shell pwd))
PROTO_DEP = ./proto_dep

GOSRC = $(shell find . -name '*.go')
BUILD_DIR = target

all: main web

proto_src = $(shell find $(PROTO_DIR) -name '*.proto')
proto_target = $(patsubst $(PROTO_DIR)/%.proto,$(PB_DIR)/%.pb.go,$(proto_src))
proto_yaml = $(patsubst $(PROTO_DIR)/%.proto,$(PB_DIR)/%.yaml,$(proto_src))

$(PB_DIR)/%.pb.go: $(PROTO_DIR)/%.proto $(PROTO_DIR)/%.yaml
	$(PROTOC) -I $(PROTO_DIR) \
		-I$(PROTO_DEP)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis\
		--go_out=plugins=grpc:$(PB_DIR) $<
	$(PROTOC) -I $(PROTO_DIR) \
		-I$(PROTO_DEP)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis\
		--grpc-gateway_out=logtostderr=true,grpc_api_configuration=$(patsubst %.proto,%.yaml,$<):. \
		--go_out=plugins=grpc:$(PB_DIR) $<
	$(PROTOC) -I $(PROTO_DIR) \
		-I$(PROTO_DEP)/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis\
		--swagger_out=logtostderr=true,grpc_api_configuration=$(patsubst %.proto,%.yaml,$<):$(PROTO_DIR)/doc \
		--go_out=plugins=grpc:$(PB_DIR) $<

proto: $(proto_target)
	cp -R $(PROJNAME)/* ./
	rm -rf $(PROJNAME)

main: $(BUILD_DIR)/$(PROJNAME)
	
web: $(BUILD_DIR)/web

$(BUILD_DIR)/$(PROJNAME): $(GOSRC) gopath
	cd $(PROJ_GOPATH) && go build -o $@ main.go

$(BUILD_DIR)/web: $(GOSRC) gopath
	cd $(PROJ_GOPATH) && go build -o $@ web/main.go

clean:
	@-rm -f $(proto_target) $(BUILD_DIR)/*

include ./make_test.mk
include make_gopath.mk