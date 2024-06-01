FROM ubuntu:22.04

RUN apt update

RUN apt-get install -y make curl wget vim

RUN wget https://github.com/YosysHQ/oss-cad-suite-build/releases/download/2024-05-31/oss-cad-suite-linux-x64-20240531.tgz
#RUN wget https://github.com/YosysHQ/oss-cad-suite-build/releases/download/2024-05-31/oss-cad-suite-linux-arm64-20240531.tgz
#RUN oss-cad-suite-darwin-arm64-20240531.tgz

RUN tar -xzvf oss-cad-suite-linux-x64-20240531.tgz

ENTRYPOINT ["bash", "-c", "source oss-cad-suite/environment && bash"]

