#!/usr/bin/env bash


export queue_config=/Users/liaozhouping/www/hb/queue/config/config.toml
pm2 start  "/home/ubuntu/www/queue/bin/queue  --config=${queue_config}" --name="queue"  --namespace="queue"