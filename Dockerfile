FROM golang:1.18.10 as BUILDER

# build binary
COPY . /go/src/github.com/VictorZhou123/xihe-forward
RUN cd /go/src/github.com/VictorZhou123/xihe-forward && GO111MODULE=on CGO_ENABLED=0 go build

# copy binary config and utils
FROM alpine:latest

RUN adduser mindspore -u 5000 -D
USER mindspore
WORKDIR /opt/app/

COPY  --from=BUILDER /go/src/github.com/VictorZhou123/xihe-forward/xihe-forward /opt/app

ENTRYPOINT ["/opt/app/xihe-forward"]
