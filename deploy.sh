#!/bin.bash

export $(grep -v '^#' deploy.env | xargs -d '\n')
docker build . -t $BUILD_TAG
docker save $BUILD_TAG:latest | ssh -C $REMOTE docker load

while true; do

read -p "Do you want to restart remote deplayment? (y/n) " yn

case $yn in 
	[yY] )
		break;;
	[nN] )
		exit;;
	* ) echo invalid response;;
esac

done

ssh -C $REMOTE docker compose -f $REMOTE_COMPOSE_FILE up -d
ssh -C $REMOTE docker compose -f $REMOTE_COMPOSE_FILE logs