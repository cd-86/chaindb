#! /bin/bash

if [ linux-gnu != $OSTYPE ]; then
    echo "该脚本只适用于 GNU+Linux"
    exit 1
fi

cd /tmp
mkdir -p shynur/Downloads
cd shynur/Downloads

# 安装 Go 工具链.
sudo rm -rf /usr/local/go
precompiled_go=go1.24.5.linux-amd64.tar.gz
if ! [ -f $precompiled_go ]; then
    wget https://go.dev/dl/$precompiled_go
fi
sudo tar -C /usr/local -xzf $precompiled_go
export PATH+=:/usr/local/go/bin
go version

# 安装 protoc.
protoc_version=30.2
if ! [ -f protoc-$protoc_version-linux-x86_64.zip ]; then
    curl -LO  \
        https://github.com/protocolbuffers/protobuf/releases/download/v$protoc_version/protoc-$protoc_version-linux-x86_64.zip
fi
mkdir -p ~/.local
unzip -o protoc-$protoc_version-linux-x86_64.zip -d $HOME/.local
export PATH+=:~/.local/bin
protoc --version

# 安装 gRPC.
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH+=:`go env GOPATH`/bin
protoc-gen-go --version


echo '~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~'
echo '请执行'
echo 'export PATH+=:/usr/local/go/bin:~/.local/bin:`/usr/local/go/bin/go env GOPATH`/bin'
