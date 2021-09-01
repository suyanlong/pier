#!/usr/bin/env bash

set -e

source x.sh

# $1 is arch, $2 is source code path
case $1 in
linux-amd64)
  print_blue "Compile for linux/amd64"
  docker run -t \
    -v $2:/code/sidercar \
    -v $2/../sidercar-client-fabric:/code/sidercar-client-fabric \
    -v $2/../sidercar-client-ethereum:/code/sidercar-client-ethereum \
    -v ~/.ssh:/root/.ssh \
    -v ~/.gitconfig:/root/.gitconfig \
    -v $GOPATH/pkg/mod:$GOPATH/pkg/mod \
    sidercar-ubuntu/compile \
    /bin/bash -c "go env -w GO111MODULE=on &&
      go env -w GOPROXY=https://goproxy.cn,direct &&
      go get -u github.com/gobuffalo/packr/packr &&
      cd /code/sidercar-client-fabric && make fabric1.4 &&
      cd /code/sidercar-client-ethereum && make eth &&
      cd /code/sidercar && make install &&
      mkdir -p /code/sidercar/bin &&
      cp /go/bin/sidercar /code/sidercar/bin/sidercar_linux-amd64 &&
      cp /code/sidercar-client-fabric/build/fabric-client-1.4.so /code/sidercar/bin/sidercar-fabric-linux.so &&
      cp /code/sidercar-client-ethereum/build/eth-client.so /code/sidercar/bin/sidercar-eth-linux.so"
  ;;
*)
  print_red "Other architectures are not supported yet"
  ;;
esac
