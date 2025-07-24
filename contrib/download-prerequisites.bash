#! /bin/bash

if [ linux-gnu != $OSTYPE ]; then
    echo "该脚本只适用于 GNU+Linux"
    exit 1
fi

# 安装 Go 工具链.
sudo rm -rf /usr/local/go
cd /tmp
mkdir -p shynur/Downloads
cd shynur/Downloads
wget https://go.dev/dl/go1.24.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz
export PATH+=:/usr/local/go/bin
go version

# 安装 gRPC 及其依赖.
