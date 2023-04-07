load('ext://ko', 'ko_build')
load('ext://deployment', 'deployment_create')

# load dev dependencies
load_dynamic('./dev/Tiltfile.dep')

# service build and deploy
ko_build('conjur-service-broker', './cmd/conjur_service_broker')

deployment_create('conjur-service-broker', 'conjur-service-broker', ports=['8080:8080'], env=read_yaml('./.env.yaml'))

k8s_resource('conjur-service-broker', port_forwards=['8080'], labels=['conjur-service-broker'], resource_deps=['api_key'])

# integration tests
load_dynamic('./test/integration/Tiltfile.dep')

# deploy ruby version of service broker for testing purposes
load_dynamic('./Tiltfile.ruby')
