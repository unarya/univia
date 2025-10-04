#!/bin/bash
cd infra && docker compose down
cd ../api/gin && docker build -t univia:dev --target=development .
cd ../../signaling && docker build -t univia-signaling:dev --target=development .
cd ../infra && docker compose up -d
cd .. && docker ps

echo "✅ Univia microservices are running...."