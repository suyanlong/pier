# sidercar

![build](https://github.com/link33/sidercar/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/meshplus/sidercar/branch/master/graph/badge.svg)](https://codecov.io/gh/meshplus/sidercar)

## Build

Using the follow command to install necessary tools.

```bash
make prepare
```

And then install sidercar using the following command.

```bash
make install
```

## Initialization

Using the follow command to initialize sidercar.
```bash
sidercar init
```
Default repo path is `~/sidercar`. If you want to specify the repo path, you can use `--repo` flag.

```bash
sidercar init --repo=$HOME/sidercar
```

After initializing sidercar, it will generate the follow directory:

```
~/sidercar
├── sidercar.toml
├── key.json

```

## Configuration

```toml
title = "sidercar"

[port]
pprof = 44555

[log]
level = "debug"
dir = "logs"
filename = "sidercar.log"
report_caller = false

[bitxhub]
addr = "localhost:60011"
validators = [
      "0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd",
      "0xe93b92f1da08f925bdee44e91e7768380ae83307",
      "0xb18c8575e3284e79b92100025a31378feb8100d6",
      "0x856E2B9A5FA82FD1B031D1FF6863864DBAC7995D",
]

[appchain]
plugin = "fabric-client-1.4.so"
config = "fabric"
```

`port.pprof`: the pprof server port

`log.level`: log level: debug, info, warn, error, fatal

`bitxhub.addr`: bitxhub grpc server port

`bitxhub.validators`: bitxhub validator's addresses

`appchain.plugin`: relative path in sidercar repo of appchain plugin

`appchain.config`: relative path of appchain config directory

## Usage

More details about usage is in [sidercar handbook](https://github.com/link33/sidercar/wiki/sidercar%E4%BD%BF%E7%94%A8%E6%96%87%E6%A1%A3)

## License

The sidercar library (i.e. all code outside of the cmd and internal directory) is licensed under the GNU Lesser General Public License v3.0, also included in our repository in the LICENSE.LESSER file.

The sidercar binaries (i.e. all code inside of the cmd and internal directory) is licensed under the GNU General Public License v3.0, also included in our repository in the LICENSE file.