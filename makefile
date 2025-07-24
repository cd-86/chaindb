SHELL = /bin/bash -O globstar

chaindb.exe: ./validator/block_cdn/service.pb.go
	go build -o chaindb.exe ./cmd

./validator/block_cdn/%.pb.go: ./validator/block_cdn/service.proto
	protoc  --go_out=. --go-grpc_out=.  \
		--go_opt=module=github.com/shynur/chaindb  \
		$<
