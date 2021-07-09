# Makefile for various commands

# Env Vars 
APP_NAME=goroach
APP_VERSION=1.1
NETWORK=roachnet

SERVICE_PORT=8880
DB_SERVER=db
DB_IMAGE=postgres:9.6.22-alpine
DB_PORT=5432
DB_USER=postgres
DB_DATABASE=testdb
DB_PASSWORD=postgres
DEFAULT_PAGE_SIZE=20

build:
	docker build -t ${APP_NAME}:${APP_VERSION} .

network:
	docker network create -d bridge ${NETWORK}

runapp:
	docker run --rm -d --net ${NETWORK} --name ${APP_NAME} \
		-e SERVICE_PORT=${SERVICE_PORT}  -e DB_SERVER=${DB_SERVER} \
		-e DB_PORT=${DB_PORT} -e DB_USER=${DB_USER} -e DB_DATABASE=${DB_DATABASE} \
		-e DB_PASSWORD=${DB_PASSWORD} -e DEFAULT_PAGE_SIZE=${DEFAULT_PAGE_SIZE} \
		-p ${SERVICE_PORT}:${SERVICE_PORT} ${APP_NAME}:${APP_VERSION}

rundb:
	docker run --rm -d --net ${NETWORK} --name $(DB_SERVER) \
		-e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} \
		-p ${DB_PORT}:${DB_PORT} $(DB_IMAGE)

clean:
	docker stop ${APP_NAME}
	docker stop ${DB_SERVER}
	docker network rm ${NETWORK}
