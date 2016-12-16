FROM golang:latest
RUN go get github.com/nats-io/gnatsd && \
    go get github.com/fanux/lhttp 
WORKDIR /go/src/github.com/fanux/lhttp
RUN git clone https://github.com/fanux/lhttp-web-demo && cd websocketServer && go install
CMD ./start.sh
