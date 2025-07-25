#! /bin/bash

if [ linux-gnu != $OSTYPE ]; then
    echo "该脚本只适用于 GNU+Linux"
    exit 1
fi

cd /tmp
mkdir -p shynur/Downloads
cd shynur/Downloads

# 安装 Go 工具链.
go_version=1.24.5
export PATH+=:/usr/local/go/bin
if ! go version | grep -F $go_version >/dev/null; then
    sudo rm -rf /usr/local/go
    precompiled_go=go$go_version.linux-amd64.tar.gz
    if ! [ -f $precompiled_go ]; then
        wget https://go.dev/dl/$precompiled_go
    fi
    sudo tar -C /usr/local -xzf $precompiled_go
fi
go version

# 安装 protoc.
protoc_version=30.2
mkdir -p ~/.local
export PATH+=:~/.local/bin
if ! protoc --version | grep -F $protoc_version >/dev/null; then
    if ! [ -f protoc-$protoc_version-linux-x86_64.zip ]; then
        curl -LO  \
            https://github.com/protocolbuffers/protobuf/releases/download/v$protoc_version/protoc-$protoc_version-linux-x86_64.zip
    fi
    unzip -o protoc-$protoc_version-linux-x86_64.zip -d $HOME/.local
fi
protoc --version

# 安装 gRPC.
export PATH+=:`go env GOPATH`/bin
if ! protoc-gen-go --version; then
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi
if ! protoc-gen-go-grpc --version; then
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi
protoc-gen-go --version
protoc-gen-go-grpc --version

echo '~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~'
echo '请执行'
echo 'export PATH+=:/usr/local/go/bin:~/.local/bin:`/usr/local/go/bin/go env GOPATH`/bin'
