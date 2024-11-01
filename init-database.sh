#!/bin/bash

echo "initializing database";
psql -U postgres <<-EOSQL
    CREATE DATABASE outagealert;
    GRANT ALL PRIVILEGES ON DATABASE outagealert TO docker;
EOSQL
