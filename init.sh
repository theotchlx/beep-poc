#!/bin/bash

# Create the messages index with the necessary mappings
curl -X PUT "elasticsearch:9200/messages" -H 'Content-Type: application/json' -d'
{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "author": { "type": "keyword" },
      "createdAt": { "type": "date" },
      "content": { "type": "text" }
    }
  }
}'

echo "Elasticsearch index 'messages' created."
