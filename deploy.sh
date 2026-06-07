#!/bin/bash
set -e

VPS_USER="trbiv"
VPS_HOST="157.230.84.210"
DEPLOY_DIR="/home/trbiv/iv-blog"

echo "Building..."
yarn build

echo "Deploying to $VPS_USER@$VPS_HOST:$DEPLOY_DIR..."
ssh "$VPS_USER@$VPS_HOST" "mkdir -p $DEPLOY_DIR"
scp -r build/* "$VPS_USER@$VPS_HOST:$DEPLOY_DIR/"

echo "Done."
