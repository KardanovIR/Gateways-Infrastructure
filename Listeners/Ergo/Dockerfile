FROM golang:alpine as builder
ENV GO111MODULE=on
ADD . /go/src/github.com/wavesplatform/GatewaysInfrastructure/Listeners/Ergo
WORKDIR /go/src/github.com/wavesplatform/GatewaysInfrastructure/Listeners/Ergo

RUN apk add build-base
RUN apk update
RUN apk add git

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine
COPY --from=0 /go/src/github.com/wavesplatform/GatewaysInfrastructure/Listeners/Ergo /
COPY config/config.y*ml /config

EXPOSE 5006
CMD ["/main"]
