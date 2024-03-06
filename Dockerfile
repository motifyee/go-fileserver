
FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN eport GO111MODULE=on
RUN go get https://github.com/motifyee/go-fileserver/master
RUN cd /build && git clone https://github.com/motifyee/go-fileserver.git

RUN cd /build/gowasm/fileserver && go build

EXPOSE 8989

ENTRYPOINT [ "/build/gowasm/fileserver" ]