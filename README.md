# beep-poc for Keycloak authn and Elasticsearch messages and search

A Proof of Concept project messages web app with Keycloak authentication and Elasticsearch full-text search.

## Launch and configuration

Local stack setup: (without service mesh)

```bash

docker compose up -d keycloak
...

# Let's configure keycloak: at port 7080
# Create realm: beep-poc
# In realm settings, goto Login tab, check User registration On.
# Create client: beep-poc-front, Valid redirect URIs as http://localhost:4040/*, Web origins as http://localhost:4040
# In clients, goto beep-poc-front, scroll a bit - under Capability config, check Direct access grants to be checked. (And Standard flow, if that isn't already the case). Uncheck OAuth 2.0 Device Authorization Grant if that isn't already the case.


docker compose up -d elasticsearch
...

# Wait a bit for elasticsearch to init.

docker compose up -d init-elasticsearch
...

# This creates the indexes, if they don't already exist in the database.

cd backend

export ES_ADDRESS=http://0.0.0.0:9200
export ES_USERNAME=elastic
export ES_PASSWORD=thisisaverystrongpassword
go run main.go
...

cd ..
cd frontend

docker compose up -d
...

```

Some commands to try things out: (see backend's readme for more tests)

```bash
# Authenticate to keycloak with curl, to retrieve access token.
curl -X POST http://localhost:7080/realms/beep-poc/protocol/openid-connect/token -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=password" -d "client_id=beep-poc-front" -d "username=abc" -d "password=abc"

# Authenticate a request:
curl -X GET "http://localhost:8080/messages/0a3ffd3b-ade0-42a9-83e7-d7fab82de051" -H "Authorization: Bearer <my access token here>"
```

Service mesh setup:

```bash

# Start by configuring the cluster

minikube start --memory=6144 --driver=kvm2 --container-runtime=containerd
```

This my Minikube configuration, cleaner than the default, especially for the service mesh installation.

```bash

# Install the Linkerd CLI on your local machine

curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install-edge | sh

export PATH=$HOME/.linkerd2/bin:$PATH

linkerd version # You should see a client version, and an unavailable server version (because we haven't installed Linkerd on the cluster yet)

# Install the Gateway API CRDs
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.1/standard-install.yaml

# Check that the cluster is ready for the control plane to be installed:
linkerd check --pre # All checks should be green.

# Install the Linkerd control plane

linkerd install --crds | kubectl apply -f -
linkerd install | kubectl apply -f -

# Wait until all installation checks are green:
linkerd check

# Deploy the applications and inject them in the Linkerd control plane simultaneously.
linkerd inject deployments/. | kubectl apply -f -
```

If you want, you can install the Linkerd visualization dashboard:

```bash
linkerd viz install | kubectl apply -f - # install the on-cluster metrics stack
linkerd check # May take a while. Everything should be green.
linkerd viz dashboard & # Access the dashboard
```
