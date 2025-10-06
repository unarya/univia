#!/bin/bash
set -e
set -o pipefail

# ==============================================================================
# UNIVIA - DEPLOY SCRIPT
# ------------------------------------------------------------------------------
# Description : Build & launch Univia microservices (API, Signaling, Infra)
# Author      : Tiecont
# Version     : 1.0
# ==============================================================================

# ==== Configuration ====
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$ROOT_DIR/infra"
API_DIR="$ROOT_DIR/cmd/api"
SIGNALING_DIR="$ROOT_DIR/cmd/signaling"

MYSQL_IMAGE_TAG="univia:mysql"
API_IMAGE_TAG="univia:dev"
SIGNALING_IMAGE_TAG="univia-signaling:dev"

# ==== Utilities ====
log() {
    echo -e "\033[1;34m[INFO]\033[0m $1"
}

error_exit() {
    echo -e "\033[1;31m[ERROR]\033[0m $1"
    exit 1
}

# ==== Preflight ====
for dir in "$INFRA_DIR" "$API_DIR" "$SIGNALING_DIR"; do
    [ -d "$dir" ] || error_exit "Missing directory: $dir"
done

cd "$ROOT_DIR"

# ==== Steps ====

log "Pulling MySQL base image..."
docker pull mysql:latest || error_exit "Failed to pull MySQL image."
docker tag mysql:latest "$MYSQL_IMAGE_TAG"

log "Stopping old containers..."
cd "$INFRA_DIR" && docker compose down || error_exit "Failed to stop old containers."
# ----------------------------------------------------------------------
# BUILD STAGE
# ----------------------------------------------------------------------
log "Building API (Gin) image from: $API_DIR"
cd "$API_DIR"
docker build \
    --target=development \
    -t "$API_IMAGE_TAG" \
    --build-arg CACHEBUST=$(date +%s) \
    . || error_exit "Failed to build API image."

log "Building Signaling (WebRTC) image from: $SIGNALING_DIR"
cd "$SIGNALING_DIR"
docker build \
    --target=development \
    -t "$SIGNALING_IMAGE_TAG" \
    --build-arg CACHEBUST=$(date +%s) \
    . || error_exit "Failed to build signaling image."

# ----------------------------------------------------------------------
# DEPLOY STAGE
# ----------------------------------------------------------------------
log "Starting infrastructure stack..."
cd "$INFRA_DIR" && docker compose --env-file ../configs/.env up -d || error_exit "Failed to start infra stack."

log "Listing running containers..."
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo ""
echo "✅ Univia microservices are running successfully!"
echo "------------------------------------------------"
echo "  • API Image:         $API_IMAGE_TAG"
echo "  • Signaling Image:   $SIGNALING_IMAGE_TAG"
echo "  • Infra Compose:     $INFRA_DIR"
echo "------------------------------------------------"
echo ""

