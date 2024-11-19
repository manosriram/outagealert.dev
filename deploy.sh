#!/bin/bash

# Add host to known hosts
ssh-keyscan -R $OUTAGEALERT_IP >> ~/.ssh/known_hosts

# Login to Docker registry
echo $DOCKER_REGISTRY_PAT | docker login -u manosriram --password-stdin

# Build and push Docker image
docker build . -t manosriram/outagealert:app
docker push manosriram/outagealert:app

# Deploy command
ssh -v root@$OUTAGEALERT_IP "
		# Docker login
		export DOPPLER_TOKEN=$DOPPLER_TOKEN
		export POSTGRES_USER=postgres
		export DOCKER_HOST=unix:///var/run/docker.sock
		echo $DOCKER_REGISTRY_PAT | docker login -u manosriram --password-stdin

		# Remove old image if exists
		yes | (docker rmi manosriram/outagealert:app 2>/dev/null || true);

		# Update code
		cd /root/dev/outagealert.io
		git pull origin main

		# Pull new image
		docker pull manosriram/outagealert:app

		docker-compose down && docker-compose pull && docker-compose --project-name outagealert up -d --force-recreate

		# Setup Doppler
		curl -Ls https://cli.doppler.com/install.sh | sh
		doppler configure set token $DOPPLER_TOKEN

		yes | docker system prune -a;
"
