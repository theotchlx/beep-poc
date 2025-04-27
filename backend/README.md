# POC Backend

## Running the backend

Execute the following command: `go run main.go` from this `README.md`'s directory.

## Trying it out

To fetch unauthenticated endpoints:

Create (or update, by message ID) a message:

```bash
$ curl -X POST 'http://localhost:8080/messages' -H "Content-Type: application/json" -d '{"author":"Johan Dome", "content":"Hallo World!"}'

{"messageId":"abe5eb64-b159-4ae1-9c8a-34d7a2d33d48"}
```

Get a message by its ID:

```bash
$ curl -X GET 'http://localhost:8080/messages/abe5eb64-b159-4ae1-9c8a-34d7a2d33d48'

{"id":"abe5eb64-b159-4ae1-9c8a-34d7a2d33d48","author":"Johan Dome","createdAt":"2025-04-27T18:11:02.20737248+02:00","content":"Hallo World!"}
```

Get paginated messages (50 first messages):

```bash
$ curl -X GET 'http://localhost:8080/messages?limit=50&offset=0'

[{"id":"abe5eb64-b159-4ae1-9c8a-34d7a2d33d48","author":"Johan Dome","createdAt":"2025-04-27T11:49:29.43003473+02:00","content":"Hallo, world!"}]
```

Search query "hello" and get the 10 first relevant messages:

```bash
$ curl -X GET 'http://localhost:8080/search/messages?query=hallo&limit=10&offset=0'

[{"id":"abe5eb64-b159-4ae1-9c8a-34d7a2d33d48","author":"Johan Dome","createdAt":"2025-04-27T11:49:29.43003473+02:00","content":"Hallo, world!"}]
```
