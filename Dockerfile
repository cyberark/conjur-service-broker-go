FROM golang:alpine as builder
RUN apk add --no-cache upx
WORKDIR /opt/conjur_service_broker
COPY go.* .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/conjur_service_broker

RUN upx --lzma /opt/conjur_service_broker/conjur_service_broker

FROM busybox

WORKDIR /opt/conjur_service_broker

COPY --from=builder /opt/conjur_service_broker/conjur_service_broker /opt/conjur_service_broker

CMD /opt/conjur_service_broker/conjur_service_broker
