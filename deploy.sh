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
	export DOPPLER_TOKEN=$DOPPLER_TOKEN && \
  echo $DOCKER_REGISTRY_PAT | docker login -u manosriram --password-stdin && \
  
  # Remove old image if exists
  yes | (docker rmi manosriram/outagealert:app 2>/dev/null || true) && \
  
  # Cleanup
  docker system prune -a && \
  docker volume prune -a && \
  
  # Update code
  cd /root/dev/outagealert.io && \
  git pull origin main && \
  
  # Pull new image
  docker pull manosriram/outagealert:app && \
  
  # Stop existing stack and remove network
  docker stack rm outagealert && \
  docker network rm public && \
  
  # Create new overlay network with subnet configuration
	docker network create --scope=swarm --attachable -d overlay public && \
  
  # Deploy stack
	docker stack config -c docker-compose.yml | docker stack deploy -c - outagealert && \
  
  # Setup Doppler
  curl -Ls https://cli.doppler.com/install.sh | sh && \
  doppler configure set token $DOPPLER_TOKEN
  
  # Start containers
  # DOPPLER_TOKEN=$DOPPLER_TOKEN docker-compose -f /root/dev/outagealert.io/docker-compose.yml up --force-recreate -d && \
  
  # Final cleanup
  # docker system prune -a && \
  # docker volume prune -a
"
