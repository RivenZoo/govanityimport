.PHONY: main gopath proto

include make_gopath.mk

PROTOC = protoc
PB_DIR = .
PROTO_DIR = api_define
PROJNAME = $(notdir $(PROJ_GOPATH))
PROTO_DEP = ./proto_dep

GOSRC = $(shell find . -name '*.go')
BUILD_DIR = target

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

$(BUILD_DIR)/$(PROJNAME): $(GOSRC) gopath
	cd $(PROJ_GOPATH) && go build -o $@ main.go

clean:
	@-rm -f $(proto_target) $(BUILD_DIR)/*
