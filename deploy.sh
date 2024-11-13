#!/bin/bash

# Add host to known hosts
ssh-keyscan $OUTAGEALERT_IP >> ~/.ssh/known_hosts

# Login to Docker registry
echo $DOCKER_REGISTRY_PAT | docker login -u manosriram --password-stdin

# Build and push Docker image
docker build . -t manosriram/outagealert:app
docker push manosriram/outagealert:app

# Deploy command
ssh -v root@$OUTAGEALERT_IP "
		# Docker login
		export DOPPLER_TOKEN=$DOPPLER_TOKEN
		echo $DOCKER_REGISTRY_PAT | docker login -u manosriram --password-stdin

		NETWORK_NAME=outagenet

		# Remove old image if exists
		yes | (docker rmi manosriram/outagealert:app 2>/dev/null || true)

		# Cleanup
		docker system prune -a
		docker volume prune -a

		# Update code
		cd /root/dev/outagealert.io
		git pull origin main

		# Pull new image
		docker pull manosriram/outagealert:app

		# Stop existing stack and remove network
		# docker stack rm outagealert

		# Wait for stack to be fully removed
		# while docker stack ls | grep -q 'outagealert'; do
				# echo 'Waiting for stack to be removed...'
				# sleep 5
		# done

		# sleep 5

		# Check if network exists
		# if ! docker network ls --format '{{.Name}}' | grep -q "^${NETWORK_NAME}$"; then
				# echo "Creating overlay network: ${NETWORK_NAME}"
				
				# # Create network with error handling
				# if ! docker network create \
						# "${NETWORK_NAME}"; then
						# echo "Error: Failed to create network ${NETWORK_NAME}"
						# exit 1
				# fi
				
				# echo "Network ${NETWORK_NAME} created successfully"
		# else
				# echo "Network ${NETWORK_NAME} already exists"
		# fi

		# Deploy stack
		# docker stack config -c docker-compose.yml | docker stack deploy -c - outagealert
		docker-compose down && docker-compose pull && docker-compose up -d --force-recreate

		# Setup Doppler
		curl -Ls https://cli.doppler.com/install.sh | sh
		doppler configure set token $DOPPLER_TOKEN
"
