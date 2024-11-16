#!/bin/bash

echo "initializing database";
psql -U postgres <<-EOSQL
    CREATE DATABASE outagealert;
		\c outagealert;
EOSQL
