#!/bin/sh
echo "Waiting for Kafka Connect to be ready..."
# Wait until Kafka Connect REST API is available
while true; do
  http_code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 2 http://host.docker.internal:8083/)
  if [ "$http_code" -eq 200 ]; then
    echo "Kafka Connect is ready 🚀"
    break
  else
    echo "Kafka Connect not ready yet (HTTP $http_code). Retrying in 3s..."
    sleep 3
  fi
done

# Check if connector already exists
if curl -s http://host.docker.internal:8083/connectors/mongo-outbox-connector | grep -q name; then
  echo "Connector 'mongo-outbox-connector' already exists. Skipping registration."
else
  echo "Registering Debezium connector..."
  response=$(curl -s -o /dev/null -w "%{http_code}" -X POST http://host.docker.internal:8083/connectors \
    -H "Content-Type: application/json" \
    -d @/kafka/connectors/mongo-outbox-connector.json)

  if [ "$response" -eq 201 ]; then
    echo "Connector registered ✅"
  elif [ "$response" -eq 409 ]; then
    echo "Connector already exists ⚠️"
  else
    echo "Failed to register connector ❌ (HTTP $response)"
  fi
fi
