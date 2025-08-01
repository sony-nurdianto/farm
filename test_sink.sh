#!/bin/bash

BROKER="localhost:29092"
TOPIC="test-insert"
SCHEMA_REGISTRY_URL="http://localhost:8081"

SCHEMA='{
  "type": "record",
  "name": "TestMessage",
  "fields": [
    {"name": "message", "type": "string"}
  ]
}'

echo "Mengirim data ke topic $TOPIC ..."

kafka-avro-console-producer \
  --bootstrap-server $BROKER \
  --topic $TOPIC \
  --property value.schema="$SCHEMA" <<EOF
{"message": "Hello Kafka via script!"}
EOF

echo "Selesai."
