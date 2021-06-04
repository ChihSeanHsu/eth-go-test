#!/bin/sh

if [ -z $1 ]; then
  docker run --network=deployment_example -d deployment_app /app/worker
else
  docker run --network=deployment_example -d deployment_app /app/worker -start $1
fi