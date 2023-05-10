ARG BUILDER_IMAGE=golang:1.20-alpine
ARG BASE_IMAGE=busybox

FROM ${BUILDER_IMAGE} as builder
RUN apk add --no-cache upx
WORKDIR /src
COPY go.* .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/conjur_service_broker

RUN upx --lzma /src/conjur_service_broker

FROM ${BASE_IMAGE}

WORKDIR /opt/conjur_service_broker

COPY --from=builder /src/conjur_service_broker /opt/conjur_service_broker

CMD /opt/conjur_service_broker/conjur_service_broker
