#!/bin/bash

# Deployment environment argument (default to dev if not provided)
ENVIRONMENT=${1:-dev}

# ECS Server SSH details
if [ "$ENVIRONMENT" = "prod" ]; then
    ECS_HOST="memowise.ink"
else
    ECS_HOST="memowise.ink"  # todo Replace with your development server
fi

ECS_USER="root"

# Check if local branch is up to date with remote branch
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse "origin/main")

if [ $LOCAL != $REMOTE ]; then
    echo "Local branch is not up to date with remote. Pushing changes..."
    git push
    if [ $? -ne 0 ]; then
        echo "git push failed. Exiting script."
        exit 1
    fi
else
    echo "Local branch is up to date with remote."
fi

IMAGE_NAME="registry.cn-qingdao.aliyuncs.com/chaosflutter/common-go"
IMAGE_TAG=$(git rev-parse --short=6 HEAD)
CONTAINER_NAME="common-go"
PORT_MAPPING="8080:8080"
VOLUME_MAPPING="/root/audio:/root/audio"

DOCKER_BUILDKIT=1 \
DOCKER_DEFAULT_PLATFORM=linux/amd64 \
docker build -t $IMAGE_NAME:$IMAGE_TAG .

if [ $? -ne 0 ]; then
    echo "build image failed"
    exit 1
fi

docker push $IMAGE_NAME:$IMAGE_TAG
if [ $? -ne 0 ]; then
    echo "Push image to registry failed"
    exit 1
fi

# Commands to execute on the ECS server
DEPLOY_COMMANDS="
  if docker ps -a | grep -q $CONTAINER_NAME; then
    echo 'stopping container: $CONTAINER_NAME' &&
    docker stop $CONTAINER_NAME &&
    docker rm $CONTAINER_NAME
  else
    echo 'container $CONTAINER_NAME does not exist or is not running.'
  fi
  docker pull $IMAGE_NAME:$IMAGE_TAG &&
  docker run -d --name $CONTAINER_NAME -p $PORT_MAPPING -v $VOLUME_MAPPING $IMAGE_NAME:$IMAGE_TAG
  sleep 2
  docker ps
"
ssh $ECS_USER@$ECS_HOST "$DEPLOY_COMMANDS"

echo "local last commit: $LOCAL"
