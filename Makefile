PROTO_SERVICE ?= unknown
PROTO_VERSION ?= v1
GEN_PATH ?= ./

.PHONY: all proto
all: proto

proto:
	protoc --proto_path=./contracts \
		-I ./contracts/third_party/protovalidate \
		--go_out=$(GEN_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_PATH) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GEN_PATH) --grpc-gateway_opt=paths=source_relative \
		./contracts/$(PROTO_SERVICE)/$(PROTO_VERSION)/$(PROTO_SERVICE).proto