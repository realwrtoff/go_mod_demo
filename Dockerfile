FROM centos:centos7

RUN yum install -y epel-release \
    && yum reinstall -y glibc-common \
    && yum install -y make \
    && yum -y install python36u python36u-pip  python36u-devel \
    && yum install -y git gcc \
    && makedir -p /app

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >> /etc/timezone

# python3
COPY requirements.txt /
RUN pip3 install --upgrade pip  --user -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com
RUN pip3 install --user -r requirements.txt -i http://mirrors.aliyun.com/pypi/simple/ --trusted-host mirrors.aliyun.com

# golang
ENV GOPROXY=https://goproxy.io
RUN curl -OL https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz && \
    tar -xzvf go1.14.1.linux-amd64.tar.gz && mv go /usr/local
ENV PATH=$PATH:/usr/local/go/bin
ENV GOROOT=/usr/local/go

# 添加代码和编译
RUN git clone https://github.com/realwrtoff/go_mod_demo.git \
    && cd go_mod_demo && make output

EXPOSE 7060

WORKDIR /go_mod_demo/output
CMD [ "bin/echo", "-c", "configs/echo.json" ]
