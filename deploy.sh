export $(grep -v '^#' deploy.env | xargs -d '\n')
docker build . -t $BUILD_TAG
docker save $BUILD_TAG:latest | ssh -C $REMOTE docker load
ssh -C $REMOTE docker compose -f $REMOTE_COMPOSE_FILE up -d
ssh -C $REMOTE docker compose -f $REMOTE_COMPOSE_FILE logs