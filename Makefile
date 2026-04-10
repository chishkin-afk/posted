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


keys:
	openssl genrsa -out ${GEN_PATH} 2048
	openssl rsa -in ${GEN_PATH} -pubout -out ${GEN_PATH}.pub

mkcertlocal:
	mkcert -install
	mkcert localhost 127.0.0.1
	mv localhost+1.pem server.crt
	mv server.crt ${GEN_PATH}
	mv localhost+1-key.pem server.key
	mv server.key ${GEN_PATH}