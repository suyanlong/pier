#!/usr/bin/env bash

set -e

source x.sh

# $1 is arch, $2 is source code path
case $1 in
linux-amd64)
  print_blue "Compile for linux/amd64"
  docker run -t \
    -v $2:/code/sidecar \
    -v $2/../sidecar-client-fabric:/code/sidecar-client-fabric \
    -v $2/../sidecar-client-ethereum:/code/sidecar-client-ethereum \
    -v ~/.ssh:/root/.ssh \
    -v ~/.gitconfig:/root/.gitconfig \
    -v $GOPATH/pkg/mod:$GOPATH/pkg/mod \
    sidecar-ubuntu/compile \
    /bin/bash -c "go env -w GO111MODULE=on &&
      go env -w GOPROXY=https://goproxy.cn,direct &&
      go get -u github.com/gobuffalo/packr/packr &&
      cd /code/sidecar-client-fabric && make fabric1.4 &&
      cd /code/sidecar-client-ethereum && make eth &&
      cd /code/sidecar && make install &&
      mkdir -p /code/sidecar/bin &&
      cp /go/bin/sidecar /code/sidecar/bin/sidecar_linux-amd64 &&
      cp /code/sidecar-client-fabric/build/fabric-client-1.4.so /code/sidecar/bin/sidecar-fabric-linux.so &&
      cp /code/sidecar-client-ethereum/build/eth-client.so /code/sidecar/bin/sidecar-eth-linux.so"
  ;;
*)
  print_red "Other architectures are not supported yet"
  ;;
esac
