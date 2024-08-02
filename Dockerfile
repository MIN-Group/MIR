FROM golang:1.18-alpine AS build

WORKDIR $GOPATH/src

# 安装git
ENV ALPINE_MIRROR=mirrors.aliyun.com
RUN echo https://mirrors.aliyun.com/alpine/v3.16/main/ > /etc/apk/repositories

RUN apk add --no-cache git gcc libpcap-dev build-base

# 克隆minlib代码
RUN git clone https://x-access-token:JTRUGsneYyk53PuYxmdt@git.sscfs.cn/pkusz-future-network-lab/common/minlib.git
# 切换minlib分支
WORKDIR $GOPATH/src/minlib
RUN git checkout parallel-mir

# 设置 GOPROXY 代理
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn

# 安装minlib依赖
RUN go mod download

# 克隆mir-go代码
WORKDIR $GOPATH/src
RUN git clone https://x-access-token:JTRUGsneYyk53PuYxmdt@git.sscfs.cn/pkusz-future-network-lab/mir/mir-go.git
# 切换mir-go分支
WORKDIR $GOPATH/src/mir-go
RUN git checkout parallel-mir



# 安装mir-go依赖
RUN go mod tidy

# 编译mir
WORKDIR $GOPATH/src/mir-go/daemon/mircmd/
RUN go install ./mir
RUN go install ./mird
RUN go install ./mirgen

# 编译mirc
WORKDIR $GOPATH/src/mir-go/daemon/mgmt
RUN go install ./mirc

# 拷贝可执行文件到镜像里
FROM alpine
ENV ALPINE_MIRROR=mirrors.aliyun.com
RUN echo https://mirrors.aliyun.com/alpine/v3.16/main/ > /etc/apk/repositories
RUN apk add --no-cache libpcap-dev bash
WORKDIR /root/

RUN mkdir -p /usr/local/etc \
           && mkdir -p /usr/local/etc/mir \
           && mkdir -p /usr/local/etc/mir/passwd

COPY --from=build /go/bin/mir /usr/local/bin/
COPY --from=build /go/bin/mird /usr/local/bin/
COPY --from=build /go/bin/mirgen /usr/local/bin/
COPY --from=build /go/bin/mirc /usr/local/bin/
COPY --from=build /go/src/mir-go/mirconf.ini .
RUN cp mirconf.ini /usr/local/etc/mir/
COPY --from=build /go/src/mir-go/defaultRoute.xml .
RUN cp defaultRoute.xml /usr/local/etc/mir/
RUN touch /tmp/mir.sock
RUN echo "success"




