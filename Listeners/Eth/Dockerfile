FROM golang:alpine as builder
ADD . /go/src/github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth
WORKDIR /go/src/github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth

RUN apk add build-base

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine
COPY --from=0 /go/src/github.com/wavesplatform/GatewaysInfrastructure/Listeners/Eth /
COPY config/config.y*ml /config

EXPOSE 5001
CMD ["/main"]
