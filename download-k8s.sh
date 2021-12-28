#!/bin/bash

#VERSION=v1.22.5
VERSION=v1.23.1

oldwd=$(pwd)
tmp=/tmp
cd ${tmp}
rm -f kubernetes.tar.gz kubernetes
url=https://github.com/kubernetes/kubernetes/releases/download/${VERSION}/kubernetes.tar.gz
if command -v curl 2>/dev/null; then
    curl -LO ${url}
elif command -v wget 2>/dev/null; then
    wget ${url}
else
    echo "error: cannot download kubernetes per curl or wget"
    exit 1
fi

set -e
set -x

tar xvf kubernetes.tar.gz
cd kubernetes
set +e
echo "y" | ./cluster/get-kube-binaries.sh
set -e

cd ${tmp}
export PATH=${PWD}/kubernetes/client/bin:${PATH}

cd ${tmp}/kubernetes/server
    tar xvf kubernetes-server-linux-amd64.tar.gz
cd ${tmp}
export PATH=${PWD}/kubernetes/server/kubernetes/server/bin:${PATH}

echo export PATH=${PATH}
cd ${oldwd}

set +e
set +x
