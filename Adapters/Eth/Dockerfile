FROM golang:alpine as builder
ADD . /go/src/github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth
WORKDIR /go/src/github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth

RUN apk add build-base
RUN apk update
RUN apk add git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.toml Gopkg.toml
COPY Gopkg.lock Gopkg.lock

RUN dep ensure --vendor-only

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine
COPY --from=0 /go/src/github.com/wavesplatform/GatewaysInfrastructure/Adapters/Eth /
COPY config/config.y*ml /config
RUN apk add --no-cache ca-certificates && update-ca-certificates 2>/dev/null || true

EXPOSE 5001
CMD ["/main"]
