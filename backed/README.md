后端 采用go-zero框架

中间件为kafka
数据库为mysql


goctl api go -api doc/api/gosmo.api -dir . doc/api/gosmo.api

wget https://mirrors.aliyun.com/golang/go1.23.8.linux-amd64.tar.gz





FROM debian:bookworm-slim

WORKDIR /tmp

RUN sed -i 's|deb.debian.org|mirrors.aliyun.com|g' /etc/apt/sources.list && \
sed -i 's|security.debian.org|mirrors.aliyun.com|g' /etc/apt/sources.list

RUN apt-get update && \
apt-get install -y --no-install-recommends \
wget \
ca-certificates \
libpcap-dev \  
&& rm -rf /var/lib/apt/lists/*

RUN wget https://mirrors.aliyun.com/golang/go1.23.8.linux-amd64.tar.gz -O go.tar.gz && \
tar -C /usr/local -xzf go.tar.gz && \
rm go.tar.gz

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
ENV CGO_ENABLED 1 

RUN mkdir -p $GOPATH/src $GOPATH/bin $GOPATH/pkg

RUN go version && \
go env

RUN apt-get purge -y wget && \
apt-get autoremove -y && \
apt-get clean

RUN useradd -m gouser && \
chown -R gouser:gouser $GOPATH

WORKDIR $GOPATH/src/app
USER gouser

CMD ["bash"]

, \
"--input-raw=:8888", \
"--input-raw-track-response", \
"--output-elasticsearch-host=http://47.94.96.190:9200", \
"--output-elasticsearch-index=gosmo", \
"--output-stdout"

docker run -d --name agent --network host   registry.cn-shanghai.aliyuncs.com/qianjisan/gosmo-agent:latest    --input-raw=:8080  --input-raw-track-response  --output-stdout --output-elasticsearch-host=http://47.94.96.190:9200 --output-elasticsearch-index=gosmo
