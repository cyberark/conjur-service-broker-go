FROM golang:alpine as builder
WORKDIR /opt/conjur_service_broker
COPY go.* .
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" ./cmd/conjur_service_broker

FROM busybox

WORKDIR /opt/conjur_service_broker

COPY --from=builder /opt/conjur_service_broker/conjur_service_broker /opt/conjur_service_broker

CMD /opt/conjur_service_broker/conjur_service_broker
