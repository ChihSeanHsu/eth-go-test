#!/bin/sh
DB_CONN_STR="host=db user=postgres password=example dbname=db port=5432 sslmode=disable TimeZone=Asia/Taipei"
docker run --network=deployment_example -e DB_CONN_STR="${DB_CONN_STR}" -it deployment_app /app/migration