FROM golang:1.10-stretch as build

WORKDIR /go/src/github.com/AliyunContainerService/gpushare-scheduler-extender
COPY . .

ADD pkg/utils/interference.json /data/interference.json

RUN go build -o /go/bin/gpushare-sche-extender cmd/*.go

FROM debian:stretch-slim

COPY --from=build /go/bin/gpushare-sche-extender /usr/bin/gpushare-sche-extender

CMD ["gpushare-sche-extender"]
