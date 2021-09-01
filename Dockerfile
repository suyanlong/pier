# use : docker build --build-arg PLUGIN_URL=https://github.com/link33/sidercar-client-fabric --build-arg PLUGIN_VERSION=v1.0.0-rc2 -t sidercar:1.0.0 .
FROM golang:1.17 as builder

RUN mkdir -p /go/src/github.com/link33/sidercar
WORKDIR /go/src/github.com/link33/sidercar

# Cache dependencies
COPY go.mod .
COPY go.sum .

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download -x

# Build real binaries
COPY . .

RUN go get -u github.com/gobuffalo/packr/packr

RUN make install

RUN cd .. && \
    git clone https://github.com/link33/sidercar-client-fabric.git && \
    cd sidercar-client-fabric && \
    git checkout v1.1.0-rc1 && \
    make fabric1.4 && \
    cp build/fabric-client-1.4 /go/bin/fabric-client-1.4

# Final image
FROM frolvlad/alpine-glibc

WORKDIR /root

# Copy over binaries from the builder
COPY --from=builder /go/bin/sidercar /usr/local/bin
COPY ./build/libwasmer.so /lib
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/lib

RUN ["sidercar", "init"]

RUN mkdir -p /root/sidercar/plugins
COPY --from=builder /go/bin/fabric-client-1.4 /root/sidercar/plugins/appchain_plugin
COPY scripts/docker_entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker_entrypoint.sh

ENV APPCHAIN_NAME=fabric-client-1.4 \
    PLUGIN_CONFIG=fabric

EXPOSE 44555 44544

# run sidercar inside container using:
# docker run -t --env APPCHAIN_NAME=fabric-client-1.4.so \
# --env PLUGIN_CONFIG=fabric \
# -v your/path/to/rule:/root/sidercar/validating.wasm \
# -v your/path/to/sidercar.toml:/root/sidercar/sidercar.toml \
# -v your/path/to/plugin_config:/root/sidercar/fabric \
# sidercar:1.0.0

ENTRYPOINT docker_entrypoint.sh, "$APPCHAIN_NAME", "$PLUGIN_CONFIG"
