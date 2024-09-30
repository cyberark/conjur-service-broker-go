kind version >/dev/null || { echo "kind cli required"; exit 1; }
export KIND_EXPERIMENTAL_DOCKER_NETWORK="kind-network"
kind create cluster

kind_cid="$(docker inspect --format="{{.Id}}" kind-control-plane)"
kind_ip="$(dirname "$(docker network inspect kind-network | yq ".[0][\"Containers\"][\"$kind_cid\"][\"IPv4Address\"]")")"
kind_port="$(dirname "$(docker port kind-control-plane)")"
kubectl config set clusters.kind-kind.server "https://$kind_ip:$kind_port"

kubectl config set-context --current --namespace=default
