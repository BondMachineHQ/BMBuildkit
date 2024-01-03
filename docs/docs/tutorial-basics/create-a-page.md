---
sidebar_position: 1
---

# Quickstart

The following is a quick example on how to get started with a demo Lattice ice40 stick FPGA and a blinky led algorithm.

## Install bmctl CLI

```bash
git clone https://github.com/BondMachineHQ/BMBuildkit
cd BMBuildkit 
go install ./cmd/bmctl
```

## Build, tag and share your FPGA algorithms

__N.B.__ you need an host
with [icepack](https://github.com/YosysHQ/icestorm), [nextpnr](https://github.com/YosysHQ/nextpnr)
and [yosys](https://github.com/YosysHQ/yosys) installed on your host. In alternative you can try with the pre-built
steps below.

```bash
bmctl build ./ -f BMFile -t <dockerhub USERNAME>/bmtest:built
```

## (alternative) Upload a pre-built firmware

```bash
bmctl build ./ -f BMFile.localfirmware -t <dockerhub USERNAME>/bmtest:pre-built
```

## Load the firmware on the board

__N.B.__ you need an host with [iceprog](https://github.com/YosysHQ/icestorm) installed.

```bash
bmctl load <dockerhub USERNAME>/bmtest:pre-built lattice/ice40/yosys
```

## Quick build with ice40 docker

```bash
docker run -e MODULE_NAME=blinky -e SYNTH_FILE=blinky.v -v $PWD/examples/blinky:/opt/source -ti dciangot/yosys bash
```
