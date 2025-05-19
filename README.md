# beep-poc

A Proof of Concept project messages web app with Keycloak authentication and Elasticsearch full-text search.

## Launch and configuration

```bash

# Let's start.


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

Some commands, also to update the backend's readme:

```bash
# Authenticate to keycloak with curl, to retrieve access token.
curl -X POST http://localhost:7080/realms/beep-poc/protocol/openid-connect/token -H "Content-Type: application/x-www-form-urlencoded" -d "grant_type=password" -d "client_id=beep-poc-front" -d "username=abc" -d "password=abc"

# Authenticate a request:
curl -X GET "http://localhost:8080/messages/0a3ffd3b-ade0-42a9-83e7-d7fab82de051" -H "Authorization: Bearer <my access token here>"
```
