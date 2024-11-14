# ./restore_psql <file_path>
cat $1 | docker exec -i docker ps -qf "name=outagealert-db-1" psql -U $POSTGRES_USER;
