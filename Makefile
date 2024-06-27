
gen-product:
	@protoc \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	./proto/product/product.proto

run-server:
	@go run server/main.go

run-client:
	@go run client/main.go
	