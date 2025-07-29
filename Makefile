SHELL = /bin/bash -O globstar


chaindb.exe: \
		always  \
		./validator/block_cdn/service.pb.go  \
		./validator/block_cdn/service_grpc.pb.go
	go build -o chaindb.exe ./cmd

./validator/block_cdn/%.pb.go ./validator/block_cdn/%_grpc.pb.go: ./validator/block_cdn/service.proto
	protoc  \
		--go_out=.           --go_opt=module=github.com/shynur/chaindb  \
		--go-grpc_out=. --go-grpc_opt=module=github.com/shynur/chaindb  \
		$<

.PHONY: always

.PHONY: clean
clean:
	rm -f ./**/?*.exe
	rm -f ./**/?*.pb.go
	rm -f ./coverage.txt
