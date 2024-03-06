
FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
# RUN go get github.com/motifyee/go-fileserver/master
RUN cd /build && git clone https://github.com/motifyee/go-fileserver.git

RUN cd /build/go-fileserver/fileserver
RUN go build

EXPOSE 8989

ENTRYPOINT [ "/build/gowasm/fileserver" ]