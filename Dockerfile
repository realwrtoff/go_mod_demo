FROM centos:centos7 as builder

RUN yum install -y make \
    && yum install -y git gcc \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" >> /etc/timezone \
    && mkdir -p /go_mod_demo/

# golang
ENV GOPROXY=https://goproxy.io
RUN curl -OL https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz && \
    tar -xzvf go1.14.1.linux-amd64.tar.gz && mv go /usr/local
ENV PATH=$PATH:/usr/local/go/bin
ENV GOROOT=/usr/local/go

# 添加代码和编译
RUN git clone https://github.com/realwrtoff/go_mod_demo.git \
    && cd go_mod_demo && git pull && make output

FROM centos:centos7
COPY --from=builder /go_mod_demo/output/ /
EXPOSE 7060

WORKDIR /go_mod_demo
CMD [ "bin/echo", "-c", "configs/echo.json" ]
