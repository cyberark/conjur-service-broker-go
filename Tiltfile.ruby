load('ext://git_resource', 'git_resource')

def deployment_generator(resource_name, image_name, namespace='default'): # returns deployment definition yaml
    envs = read_yaml('./.env.yaml')
    return blob("""apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
  namespace: %s
  labels:
    app: %s
spec:
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - name: %s
        image: %s
        command: ['rails', 's', '-p', '3000', '-b', '0.0.0.0', '-e', 'development']
        ports:
        - containerPort: 3000
          name: http
        env:
          %r
""" % (resource_name, namespace, resource_name, resource_name, resource_name, resource_name, image_name, envs))

git_resource('conjur-service-broker-ruby', 'git@github.com:cyberark/conjur-service-broker.git#main',
    resource_deps=['api_key'], port_forwards=['3000:3000'], deployment_callback=deployment_generator)
