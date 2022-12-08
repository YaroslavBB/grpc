export ROOT=..
export GRPC=$ROOT/cmd
export DOCKER_PATH=$ROOT/docker
export IMAGES_DIR="../static/images/"

cd $DOCKER_PATH && docker compose up -d
cd $GRPC && go run .
