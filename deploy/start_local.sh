#!/usr/bin/env bash

pm2 start  "/home/ubuntu/www/queue/bin/queue  --config=/home/ubuntu/www/queue/conf/config.toml" --name="queue"  --namespace="queue"


go run main.go --config=/Users/liaozhouping/www/hb/queue/config/config.toml