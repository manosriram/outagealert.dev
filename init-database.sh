#!/bin/bash
set -e

echo "initializing database";
psql -v ON_ERROR_STOP=1 -U postgres <<-EOSQL
    CREATE DATABASE outagealert;
    GRANT ALL PRIVILEGES ON DATABASE outagealert TO docker;
EOSQL
