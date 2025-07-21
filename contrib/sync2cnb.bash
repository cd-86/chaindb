#! /bin/bash

cd `dirname "$0"`/..

rm -rf /tmp/shynur/chaindb.git
mkdir -p /tmp/shynur
cp -r . /tmp/shynur/chaindb.git

cd /tmp/shynur/chaindb.git
git remote set-url origin https://cnb.cool/shynur/chaindb.git
git push
