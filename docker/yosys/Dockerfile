FROM ubuntu:22.04 as build
LABEL authors="dciangot"

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
    pkg-config build-essential curl ca-certificates python3 python3-dev git clang bison flex libreadline-dev \
                                                                                                 gawk tcl-dev libffi-dev git mercurial graphviz   \
                                                                                                 xdot pkg-config python3 libftdi-dev \
                                                                                                python3-dev libboost-all-dev cmake libeigen3-dev

RUN git clone https://github.com/YosysHQ/yosys.git /tmp/yosys \
    && cd /tmp/yosys \
    && make -j$(nproc) \
    && make \
    && make install

RUN apt-get update -qq \
 && DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
    libftdi1-dev

ENV LDFLAGS "-Wl,--copy-dt-needed-entries"

RUN git clone https://github.com/YosysHQ/icestorm.git /tmp/icestorm \
    && cd /tmp/icestorm \
    && make -j$(nproc) \
    && make install

RUN git clone --recursive https://github.com/YosysHQ/nextpnr.git /tmp/nextpnr \
 && mkdir /tmp/nextpnr/build/ \
 && cd /tmp/nextpnr/build \
 && cmake .. \
   -DARCH=ice40 \
   -DBUILD_GUI=OFF \
   -DBUILD_PYTHON=ON \
   -DUSE_OPENMP=ON \
 && make \
 && make install


ENV OUTFILE=firmware.bin
ENV SRC=/opt/source
ENV SYNTH_FILE=synth.v
ENV MODULE_NAME=ice40
ENV BOARD_MODEL=ice40

COPY Makefile /opt/Makefile

WORKDIR /opt

CMD ["make"]
