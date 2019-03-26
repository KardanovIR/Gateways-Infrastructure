FROM golang:alpine as builder
ADD . /go/src/github.com/wavesplatform/GatewaysInfrastructure/Router
WORKDIR /go/src/github.com/wavesplatform/GatewaysInfrastructure/Router

RUN apk add build-base
RUN apk update
RUN apk add git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.toml Gopkg.toml
COPY Gopkg.lock Gopkg.lock

RUN dep ensure --vendor-only

RUN CGO_ENABLED=1 GOOS=linux go build -o main .

FROM alpine
COPY --from=0 /go/src/github.com/wavesplatform/GatewaysInfrastructure/Router /
COPY config/config.y*ml /config

EXPOSE 5001
CMD ["/main"]
