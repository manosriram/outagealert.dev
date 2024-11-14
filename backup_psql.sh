DATE=$(date +%Y-%m-%d_%H%M%S)
BACKUP_PATH=/mnt/volume_sgp1_02/psql_backups
DUMPFILE_NAME=dump_$DATE
POSTGRES_CONTAINER_ID=$(docker ps -qf "name=outagealert-db-1")

mkdir -p $BACKUP_PATH
docker exec -it $POSTGRES_CONTAINER_ID pg_dumpall -c -U postgres > $BACKUP_PATH/$DUMPFILE_NAME.sql;
