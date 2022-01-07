SAY='Hello, World!'

generate-certs:
	cd certs && \
	cfssl genkey -initca ca.json | cfssljson -bare ca && \
	cfssl gencert -ca ca.pem -ca-key ca-key.pem ca.json | cfssljson -bare server && \
	cfssl gencert -ca ca.pem -ca-key ca-key.pem client-echo-scope.json | cfssljson -bare client-echo-scope && \
	cfssl gencert -ca ca.pem -ca-key ca-key.pem client-no-scope.json | cfssljson -bare client-no-scope

test:
	go test ./...

build-protos:
	protoc \
		-I ./protos \
		--go_out=./protos --go_opt=paths=source_relative \
		--go-grpc_out=./protos --go-grpc_opt=paths=source_relative \
		./protos/api/v1/service.proto

build: build-protos
	go build -o ./bin/mtls-example main.go

run: build
	./bin/mtls-example

echo:
	grpcurl -cacert certs/ca.pem -cert certs/client-echo-scope.pem -key certs/client-echo-scope-key.pem -import-path ./protos/api/v1 -proto service.proto -d "{\"say\": \"${SAY}\"}" 127.0.0.1:9090 protos.api.v1.Service/Echo

echo-no-auth:
	grpcurl -cacert certs/ca.pem -cert certs/client-no-scope.pem -key certs/client-no-scope-key.pem -import-path ./protos/api/v1 -proto service.proto -d "{\"say\": \"${SAY}\"}" 127.0.0.1:9090 protos.api.v1.Service/Echo

echo-no-tls:
	grpcurl -insecure -import-path ./protos/api/v1 -proto service.proto -d "{\"say\": \"${SAY}\"}" 127.0.0.1:9090 protos.api.v1.Service/Echo