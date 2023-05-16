protoc -I ./pkg/api/   --go_out=./pkg/api/ --go-grpc_out=./pkg/api/ converter.proto
protoc -I ./pkg/api/   --go_out=./pkg/api/ --go-grpc_out=./pkg/api/ transfer.proto