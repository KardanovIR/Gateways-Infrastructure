FROM golang:alpine as builder
ENV GO111MODULE=on
ADD . /go/src/github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc
WORKDIR /go/src/github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc

RUN apk add build-base
RUN apk update
RUN apk add git

COPY ./go.mod ./go.sum ./
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine
COPY --from=0 /go/src/github.com/wavesplatform/GatewaysInfrastructure/Adapters/Btc /
COPY config/config.y*ml /config
RUN apk add --no-cache ca-certificates && update-ca-certificates 2>/dev/null || true

EXPOSE 5001
CMD ["/main"]
