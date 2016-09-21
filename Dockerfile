FROM golang:onbuild

RUN /usr/sbin/useradd --uid 12321 -U -m -d /go -s /bin/bash -c "Golang User" gobot
RUN chown -R gobot:gobot /go
USER gobot
