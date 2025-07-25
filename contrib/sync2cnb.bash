#! /bin/bash

cd `dirname "$0"`/..

rm -rf /tmp/shynur/chaindb.git
mkdir -p /tmp/shynur
cp -r . /tmp/shynur/chaindb.git

cd /tmp/shynur/chaindb.git
ln -s ./docs/ReadMe.md ReadMe.md
git add ReadMe.md
git commit -m '为 CNB 添加仓库及 ReadMe 文件'
git remote set-url origin https://cnb.cool/shynur/chaindb.git
git push -f
