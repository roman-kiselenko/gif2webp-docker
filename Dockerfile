FROM ubuntu:latest

LABEL version="1.0"
LABEL maintainer="shindu666@gmail.com"

RUN apt-get -y update \
 && apt-get install -y git curl wget libpng-dev libtool libgif-dev autoconf automake make gcc g++

RUN curl -O https://storage.googleapis.com/golang/go1.11.1.linux-amd64.tar.gz \
  && tar xvf go1.11.1.linux-amd64.tar.gz

ENV PATH $PATH:/go/bin

WORKDIR /usr/local/webp
RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-1.0.0.tar.gz \
      && tar -xvzf libwebp-1.0.0.tar.gz \
      && mv libwebp-1.0.0 libwebp && \
      rm libwebp-1.0.0.tar.gz && \
      cd libwebp && \
      ./configure --enable-everything && \
      make && \
      make install && \
      cd .. && \
      rm -rf libwebp

ENV PATH $PATH:/usr/local/webp/libwebp-0.5.0-linux-x86-64/bin

RUN ldconfig

RUN mkdir -p /examples

COPY test_images/* /examples/

RUN mkdir -p /gif2web_api

EXPOSE 8080 8080

COPY go.mod /gif2web_api
COPY main.go /gif2web_api
WORKDIR /gif2web_api

CMD ["go", "run", "main.go"]
