# How to run it

1. `docker-compose -f deployment/docker-compose.yaml up -d` to start service and db
2. `docker-compose -f deployment/docker-compose.yaml logs -f` 
   wait until `database system is ready to accept connections` this message appear in db_1 twice.
3. `sh script/run_migration.sh` to run db migration
4. There are two ways to run sync block worker. 
   (Because it's will run in background and not with docker-compose, you need to kill it by yourself.)
    1. `sh script/run_worker.sh {blockNum want to start}` it will run start from the blockNum you specify until you kill this process.
    2. `sh script/run_worker.sh` it will run start from the latest blockNum until you kill this process.
   
5. After all the steps above are done, you can `curl localhost:8080` with these APIs
   1. `/blocks` it will return last 10 blocks. If you give the limit (`/blocks?limit=20`), it will return the last {limit} blocks
   2. `/blocks/{blockNum}` it will return the specific block.
   3. `/transaction/{txHash}` it will return the specific transaction.
