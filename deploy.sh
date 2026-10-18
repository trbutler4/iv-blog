#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
VPS_USER="trbiv"
VPS_HOST="157.230.84.210"
DEPLOY_DIR="/home/trbiv/iv-blog"
NEXT_DIR="${DEPLOY_DIR}.next"
CADDY_CONFIG="$SCRIPT_DIR/Caddyfile"
REMOTE_CADDY_CONFIG="/home/trbiv/Caddyfile"
REMOTE_APP_DIR="/home/trbiv/.local/lib/iv-blog"

echo "Checking private-page service support..."
ssh "$VPS_USER@$VPS_HOST" 'systemctl --user show-environment >/dev/null'
LINGER="$(ssh "$VPS_USER@$VPS_HOST" "loginctl show-user '$VPS_USER' --property=Linger --value")"
if [[ "$LINGER" != "yes" ]]; then
  echo "Enable persistent user services on the VPS first: loginctl enable-linger $VPS_USER" >&2
  exit 1
fi
REMOTE_ARCH="$(ssh "$VPS_USER@$VPS_HOST" uname -m)"
case "$REMOTE_ARCH" in
  x86_64) GOARCH=amd64 ;;
  aarch64) GOARCH=arm64 ;;
  *) echo "Unsupported VPS architecture: $REMOTE_ARCH" >&2; exit 1 ;;
esac
ARTIFACT_DIR="$(mktemp -d)"
trap 'rm -rf "$ARTIFACT_DIR"' EXIT

echo "Building..."
(
	cd "$SCRIPT_DIR"
	go run ./cmd/site
	CGO_ENABLED=0 GOOS=linux GOARCH="$GOARCH" go build -o "$ARTIFACT_DIR/site" ./cmd/site
)

echo "Installing private-page service..."
ssh "$VPS_USER@$VPS_HOST" "mkdir -p '$REMOTE_APP_DIR' /home/trbiv/.config/systemd/user"
scp "$ARTIFACT_DIR/site" "$VPS_USER@$VPS_HOST:$REMOTE_APP_DIR/site.next"
scp "$SCRIPT_DIR/iv-blog.service" "$VPS_USER@$VPS_HOST:/home/trbiv/.config/systemd/user/iv-blog.service"
ssh "$VPS_USER@$VPS_HOST" \
	"chmod 700 '$REMOTE_APP_DIR/site.next' && mv '$REMOTE_APP_DIR/site.next' '$REMOTE_APP_DIR/site' && systemctl --user daemon-reload && systemctl --user enable iv-blog.service && systemctl --user restart iv-blog.service && systemctl --user is-active --quiet iv-blog.service"
LOGIN_STATUS="$(ssh "$VPS_USER@$VPS_HOST" 'curl --silent --show-error --retry 5 --retry-connrefused --retry-delay 1 --output /dev/null --write-out "%{http_code}" http://127.0.0.1:8081/megan/')"
if [[ "$LOGIN_STATUS" != "401" ]]; then
  echo "Private-page service did not return the login page (HTTP $LOGIN_STATUS)." >&2
  exit 1
fi

# Install route protection before publishing private content.
echo "Reloading Caddy..."
scp "$CADDY_CONFIG" "$VPS_USER@$VPS_HOST:$REMOTE_CADDY_CONFIG"
ssh "$VPS_USER@$VPS_HOST" \
	"caddy validate --config '$REMOTE_CADDY_CONFIG' --adapter caddyfile && caddy reload --config '$REMOTE_CADDY_CONFIG' --adapter caddyfile"

echo "Deploying to $VPS_USER@$VPS_HOST:$DEPLOY_DIR..."
tar -C "$SCRIPT_DIR/build" -czf - . | ssh "$VPS_USER@$VPS_HOST" \
	"rm -rf '$NEXT_DIR' && mkdir -p '$NEXT_DIR' && tar -xzf - -C '$NEXT_DIR' && rm -rf '$DEPLOY_DIR' && mv '$NEXT_DIR' '$DEPLOY_DIR'"

echo "Done."
