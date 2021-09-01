# sidercar

![build](https://github.com/link33/sidercar/workflows/build/badge.svg)
[![codecov](https://codecov.io/gh/link33/sidercar/branch/master/graph/badge.svg)](https://codecov.io/gh/link33/sidercar)

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

[appchain]
plugin = "fabric-client-1.4.so"
config = "fabric"
```

`port.pprof`: the pprof server port

`log.level`: log level: debug, info, warn, error, fatal

`appchain.plugin`: relative path in sidercar repo of appchain plugin

`appchain.config`: relative path of appchain config directory

