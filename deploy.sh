#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
VPS_USER="trbiv"
VPS_HOST="157.230.84.210"
DEPLOY_DIR="/home/trbiv/iv-blog"
NEXT_DIR="${DEPLOY_DIR}.next"
CADDY_CONFIG="$SCRIPT_DIR/Caddyfile"
REMOTE_CADDY_CONFIG="/home/trbiv/Caddyfile"

echo "Building..."
(
	cd "$SCRIPT_DIR"
	go run ./cmd/site
)

echo "Deploying to $VPS_USER@$VPS_HOST:$DEPLOY_DIR..."
tar -C "$SCRIPT_DIR/build" -czf - . | ssh "$VPS_USER@$VPS_HOST" \
	"rm -rf '$NEXT_DIR' && mkdir -p '$NEXT_DIR' && tar -xzf - -C '$NEXT_DIR' && rm -rf '$DEPLOY_DIR' && mv '$NEXT_DIR' '$DEPLOY_DIR'"

echo "Reloading Caddy..."
scp "$CADDY_CONFIG" "$VPS_USER@$VPS_HOST:$REMOTE_CADDY_CONFIG"
ssh "$VPS_USER@$VPS_HOST" \
	"caddy validate --config '$REMOTE_CADDY_CONFIG' --adapter caddyfile && caddy reload --config '$REMOTE_CADDY_CONFIG' --adapter caddyfile"

echo "Done."
