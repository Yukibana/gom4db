FROM golang:latest

# for Chinese developer
COPY sources.list /etc/apt

# update
RUN apt-get update -y

RUN apt-get -y install libz-dev libbz2-dev libsnappy-dev zlib1g-dev  build-essential checkinstall libgflags-dev && \
apt-get -y -t stretch-backports install librocksdb-dev

ENV CGO_LDFLAGS "-L/usr/local/lib -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy"
ENV CGO_CFLAGS  "-I/usr/local/include"