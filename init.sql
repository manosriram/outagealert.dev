SELECT 'CREATE DATABASE outagealert' 
WHERE NOT EXISTS (
    SELECT FROM pg_database WHERE datname = 'outagealert'
)\gexec
