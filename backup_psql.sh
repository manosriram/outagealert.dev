mkdir -p /mnt/volume_sgp1_02/psql_backups/
docker exec -t postgres pg_dumpall -c -U postgres > /mnt/volume_sgp1_02/psql_backups/dump_`date +%Y-%m-%d"_"%H_%M_%S`.sql;
