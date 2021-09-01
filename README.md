# sidecar

![build](https://github.com/link33/sidecar/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/link33/sidecar/branch/master/graph/badge.svg)](https://codecov.io/gh/link33/sidecar)

## Build

Using the follow command to install necessary tools.

```bash
make prepare
```

And then install sidecar using the following command.

```bash
make install
```

## Initialization

Using the follow command to initialize sidecar.
```bash
sidecar init
```
Default repo path is `~/sidecar`. If you want to specify the repo path, you can use `--repo` flag.

```bash
sidecar init --repo=$HOME/sidecar
```

After initializing sidecar, it will generate the follow directory:

```
~/sidecar
├── sidecar.toml
├── key.json

```

## Configuration

```toml
title = "sidecar"

[port]
pprof = 44555

[log]
level = "debug"
dir = "logs"
filename = "sidecar.log"
report_caller = false

[appchain]
plugin = "fabric-client-1.4.so"
config = "fabric"
```

`port.pprof`: the pprof server port

`log.level`: log level: debug, info, warn, error, fatal

`appchain.plugin`: relative path in sidecar repo of appchain plugin

`appchain.config`: relative path of appchain config directory

