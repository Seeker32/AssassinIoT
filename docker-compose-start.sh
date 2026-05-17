#!/usr/bin/env bash

set -euo pipefail

sudo chown -R 1000:1000 services/emqx
sudo docker compose up -d
