FROM ubuntu:18.04

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && \
    ln -fs /usr/share/zoneinfo/Europe/Moscow /etc/localtime

LABEL maintainer="exitstop@list.ru"

ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /tmp/zoneinfo.zip
ENV ZONEINFO /tmp/zoneinfo.zip

RUN set -eux; \
        apt-get update; \
        apt-get install -y --no-install-recommends \
        g++ \
        gcc \
        libc6-dev \
        make \
        pkg-config \
        wget \
        git ca-certificates \
        && update-ca-certificates \
        ; \
        rm -rf /var/lib/apt/lists/*

ENV PATH /usr/local/go/bin:$PATH

ENV GOLANG_VERSION 1.19

ADD https://dl.google.com/go/go1.19.linux-amd64.tar.gz /tmp/go.tar.gz

RUN tar -C /usr/local -xzf /tmp/go.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

#RUN  apt update && apt install -y git ca-certificates && update-ca-certificates
#RUN apt update && apt install -y ttf-mscorefonts-installer libicu66
WORKDIR $GOPATH/github.com/exitstop/playwright

RUN go mod init github.com/exitstop/playwright && go get github.com/playwright-community/playwright-go && \
go install github.com/playwright-community/playwright-go/cmd/playwright && \
playwright install --with-deps


WORKDIR /app
