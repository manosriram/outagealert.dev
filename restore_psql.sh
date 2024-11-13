cat $1 | docker exec -i $2 psql -U $POSTGRES_USER;
