#load('ext://tests/golang', 'test_go')
load('ext://deployment', 'deployment_create')
load('ext://uibutton', 'cmd_button', 'location', 'text_input')

# load dev dependencies
load_dynamic('./dev/Tiltfile.dep')

# service build and deploy
docker_build('conjur-service-broker', '.', dockerfile_contents="""
FROM golang:alpine as builder
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

FROM busybox
WORKDIR /opt/conjur_service_broker
COPY --from=builder /src/conjur_service_broker /opt/conjur_service_broker
CMD ["/opt/conjur_service_broker/conjur_service_broker"]
""")

deployment_create('conjur-service-broker', 'conjur-service-broker', ports=['8080:8080'], env=read_yaml('./.env.yaml'))

k8s_resource('conjur-service-broker', port_forwards=['8080'], labels=['conjur-service-broker'], resource_deps=['api_key'])

# integration tests
load_dynamic('./test/integration/Tiltfile.dep')

# tests
#test_go('tests', './...', '.', timeout='30s', extra_args=['-cover'], labels=['conjur-service-broker'])
local_resource(name='tests', cmd='./scripts/test.sh', labels=['conjur-service-broker'])

cmd_button(name='coverage report',
           resource='tests',
           argv=['go', 'tool', 'cover', '-html', 'coverage/all_no_gen'],
           text='HTML report',
           icon_name='html')
