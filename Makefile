.PHONY:

# ==============================================================================
# Local env

local:
	echo "Starting local environment"
	make jaeger
	docker-compose -f docker-compose.local.yml up --build

jaeger:
	echo "Starting jaeger containers"
	docker run -d --name jaeger \
      -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
      -p 5775:5775/udp \
      -p 6831:6831/udp \
      -p 6832:6832/udp \
      -p 5778:5778 \
      -p 16686:16686 \
      -p 14268:14268 \
      -p 14250:14250 \
      -p 9411:9411 \
      jaegertracing/all-in-one:1.21

make_cert:
	echo "Generating SSL certificates"
	sh ./user/ssl/instructions.sh

# ==============================================================================
# Linter

linter:
	echo "Starting linters"
	cd main && golangci-lint run ./...
	cd ..
	cd user && golangci-lint run ./...
	cd ..
	cd sessions && golangci-lint run ./...
	cd ..




# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)


# ==============================================================================
# Go migrate postgresql

user_dbname = user_db
user_port = 5433
user_SSL_MODE = disable

force_user_db:
	migrate -database postgres://postgres:postgres@localhost:$(user_port)/$(user_dbname)?sslmode=$(user_SSL_MODE) -path user/migrations force 1

version_user_db:
	migrate -database postgres://postgres:postgres@localhost:$(user_port)/$(user_dbname)?sslmode=$(user_SSL_MODE) -path user/migrations version

migrate_user_db_up:
	migrate -database postgres://postgres:postgres@localhost:$(user_port)/$(user_dbname)?sslmode=$(user_SSL_MODE) -path user/migrations up 1

migrate_user_db_down:
	migrate -database postgres://postgres:postgres@localhost:$(user_port)/$(user_dbname)?sslmode=$(user_SSL_MODE) -path user/migrations down 1


# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache


swagger_user_api:
	echo "Starting swagger generating"
	cd ./user && echo `pwd`
	cd ./user && swag init -g **/**/*.go