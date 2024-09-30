ARG BUILDER_IMAGE_VERSION=1.23-alpine
ARG BASE_IMAGE_VERSION=1.36.0

FROM golang:${BUILDER_IMAGE_VERSION} as builder
RUN apk add --no-cache upx
WORKDIR /src
COPY go.* .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/conjur_service_broker

RUN upx --lzma /src/conjur_service_broker

FROM busybox:${BASE_IMAGE_VERSION}

WORKDIR /opt/conjur_service_broker

COPY --from=builder /src/conjur_service_broker /opt/conjur_service_broker

CMD /opt/conjur_service_broker/conjur_service_broker
