#FROM golang:latest
FROM dev.reg.iflytek.com/base/golang:1.8.3
WORKDIR /go/src/github.com/fanux
RUN go get github.com/nats-io/gnatsd && \
    go get github.com/fanux/lhttp && \
    git clone https://github.com/fanux/lhttp-web-demo && \
    cd lhttp-web-demo && go install && \
    cd ../lhttp/websocketServer && go install
CMD sh lhttp/start.sh
