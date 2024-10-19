#!bin/bash
ssh root@OUTAGEALERT_IP "cd ~/dev/outagealert && git pull origin main && docker-compose up -d --build"
