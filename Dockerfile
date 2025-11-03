ARG BUILDER_IMAGE_VERSION=1.25-alpine
ARG BASE_IMAGE_VERSION=1.36.1

FROM golang:${BUILDER_IMAGE_VERSION} as builder
RUN apk add --no-cache upx

# On CyberArk dev laptops, golang dependencies are downloaded
# with a corporate proxy in the middle. For these connections
# succeed we need to configure the proxy CA certificate in the
# build container.
#
# To also allow this script to work on non-CyberArk laptops
# we copy the certificate into the Docker image as a (potentially
# empty) directory, rather than rely on the CA file itself.
ADD build_ca_certificate /usr/local/share/ca-certificates/
RUN update-ca-certificates

WORKDIR /src
COPY go.* .
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" ./cmd/conjur_service_broker

RUN upx --lzma /src/conjur_service_broker

FROM busybox:${BASE_IMAGE_VERSION}

WORKDIR /opt/conjur_service_broker

COPY --from=builder /src/conjur_service_broker /opt/conjur_service_broker

CMD ["/opt/conjur_service_broker/conjur_service_broker"]
