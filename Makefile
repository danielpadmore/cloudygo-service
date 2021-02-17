CONTAINER_NAME=cloudygo-service
CONTAINER_VERSION=v0.0.12
DB_CONTAINER_NAME=cloudygo-db
DB_CONTAINER_VERSION=v0.0.12

build_db:
	cd data && docker build -t ${DB_CONTAINER_NAME}:${CONTAINER_VERSION} .

build_docker: build_linux
	docker build -t ${CONTAINER_NAME}:${CONTAINER_VERSION} .

run_db:
	docker run --name cloudygo-db -e POSTGRES_DB=cloudygo -d -p 5432:5432 -v ${HOME}/docker/volumes/postgres/cloudygo1:/var/lib/postgresql/data ${DB_CONTAINER_NAME}:${DB_CONTAINER_VERSION}

run_service:
	go run .
