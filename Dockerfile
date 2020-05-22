FROM centos:centos7

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
ADD cmd /go_mod_demo/cmd/
ADD configs /go_mod_demo/configs/
ADD internal /go_mod_demo/internal/
ADD scripts /go_mod_demo/scripts/
ADD go.mod /go_mod_demo/
ADD go.sum /go_mod_demo/
ADD Makefile /go_mod_demo/
RUN cd go_mod_demo && make output

EXPOSE 7060

WORKDIR /go_mod_demo/output/tpl-go-http
CMD [ "bin/echo", "-c", "configs/echo.json" ]
