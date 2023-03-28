load('ext://ko', 'ko_build')
load('ext://deployment', 'deployment_create')

load_dynamic('./dev/Tiltfile.dep')

ko_build('conjur-service-broker', '.')

deployment_create('conjur-service-broker', 'conjur-service-broker', env=read_yaml('./.env.yaml'))

k8s_resource('conjur-service-broker', labels=['conjur-service-broker'], resource_deps=['conjur'])
