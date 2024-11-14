# ./restore_psql <file_path>
POSTGRES_CONTAINER_ID=$(docker ps -qf "name=outagealert-db-1")
docker exec -i $POSTGRES_CONTAINER_ID psql -U $POSTGRES_USER < $1
