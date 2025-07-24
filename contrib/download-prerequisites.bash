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
wget https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz
export PATH+=:/usr/local/go/bin
go version

# 安装 protoc.
curl -LO  \
    https://github.com/protocolbuffers/protobuf/releases/download/v30.2/protoc-30.2-linux-x86_64.zip
mkdir -p ~/.local
unzip protoc-30.2-linux-x86_64.zip -d $HOME/.local
export PATH+=:~/.local/bin
protoc --version
