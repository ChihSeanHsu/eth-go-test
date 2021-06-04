#!/bin/sh
DB_CONN_STR="host=db user=postgres password=example dbname=db port=5432 sslmode=disable TimeZone=Asia/Taipei"
if [ -z $1 ]; then
  docker run --network=deployment_example -e DB_CONN_STR="${DB_CONN_STR}" -d deployment_app /app/worker
else
  docker run --network=deployment_example  -e DB_CONN_STR="${DB_CONN_STR}" -d deployment_app /app/worker -start $1
fi