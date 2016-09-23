FROM alpine:latest

RUN /usr/sbin/adduser -u 12321 -h /go -s /sbin/nologin -g "Golang User" -D gobot
WORKDIR /go
USER gobot

ARG app
ENV EXECUTABLE=${app:-app}-amd64

ADD build/$EXECUTABLE /go/
ADD build/upload.gtpl /go/
CMD ["sh", "-c", "./$EXECUTABLE"]
