# this image can be built from the gorocks-docker directory
FROM gorocksdb

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# memberlist port
EXPOSE 7946
# gnet port
EXPOSE 12306
# gonet port
EXPOSE 12347
# http port
EXPOSE 12345

WORKDIR /root/go

COPY . .

RUN go build  -mod=vendor .

ENTRYPOINT ["./gom4db"]

