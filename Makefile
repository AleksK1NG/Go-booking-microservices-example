.PHONY:

local:
	echo "Starting local environment"
	docker-compose -f docker-compose.local.yml up --build

linter:
	echo "Starting linters"
	cd main && echo 'cool ' && golangci-lint run ./...
	cd ..
	cd user && golangci-lint run ./...
	cd ..
	cd sessions && golangci-lint run ./...
	cd ..

jaeger:
	echo "Starting jaeger containers"
	docker run --name jaeger \
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

force_user_db:
	migrate -database postgres://postgres:postgres@localhost:$(user_port)/$(user_dbname)?sslmode=disable -path user/migrations force 1

version_user_db:
	migrate -database postgres://postgres:postgres@localhost:$(user_port)/$(user_dbname)?sslmode=disable -path user/migrations version

migrate_user_db_up:
	migrate -database postgres://postgres:postgres@localhost:$(user_port)/$(user_dbname)?sslmode=disable -path user/migrations up 1

migrate_user_db_down:
	migrate -database postgres://postgres:postgres@localhost:$(user_port)/$(user_dbname)?sslmode=disable -path user/migrations down 1