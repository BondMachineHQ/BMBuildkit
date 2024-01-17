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

### Upload demo alveo u50 firmware

```bash
wget http://bondmachine.fisica.unipg.it/firmwares/bmfloat16.xclbin

bmctl build ./ -f BMFile.localfirmware-alveo -t <dockerhub USERNAME>/bmtest-alveo:pre-built
```

## Load the firmware on the board

__N.B.__ you need an host with [iceprog](https://github.com/YosysHQ/icestorm) installed.

```bash
bmctl load <dockerhub USERNAME>/bmtest:pre-built lattice/ice40/yosys
```

### Alveo u50 load firmware

__N.B.__ you need [Xilinx runtime](https://github.com/Xilinx/Xilinx_Base_Runtime) installed and sourced for this to run.

```bash
bmctl load  <dockerhub USERNAME>/bmtest-alveo:pre-built  xilinx/alveou50/xrt
```

## Quick build with ice40 docker

```bash
docker run \
  -e MODULE_NAME=blinky \
  -e SYNTH_FILE=blinky.v \
  -v $PWD/examples/blinky/ice40:/opt/source \
  dciangot/yosys:ice40
```

```bash
docker run \
  -e MODULE_NAME=top \
  -e SYNTH_FILE=blinky.v \
  -v $PWD/examples/blinky/primer20k:/opt/source \
  dciangot/yosys:primer20k
```
