PROTO_SERVICE ?= unknown
PROTO_VERSION ?= v1
GEN_PATH ?= ./

.PHONY: all proto mkcertlocal mkcertclient keys
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

CERT_HOST ?= 127.0.0.1

mkcertlocal:
	mkcert -install
	mkcert localhost $(CERT_HOST)
	mv localhost+1.pem server.crt
	mv server.crt ${GEN_PATH}
	mv localhost+1-key.pem server.key
	mv server.key ${GEN_PATH}

mkcertclient:
	mkcert -client -cert-file $(GEN_PATH)/client.crt -key-file $(GEN_PATH)/client.key client

quick:
	echo "generating jwt keys..."
	mkdir -p keys/
	make keys GEN_PATH=./keys/jwtRS256.key

	echo "set environment variables..."
	cp .env.example .env

	docker-compose up -d