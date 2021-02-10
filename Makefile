.PHONY:

# ==============================================================================
# Start local dev environment

local:
	echo "Starting local environment"
	#make jaeger
	docker-compose -f docker-compose.local.yml up --build

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



# ==============================================================================
# Generate swagger documentation

swagger_user_api:
	echo "Starting swagger generating"
	cd ./user && swag init -g **/**/*.go
	cd ..
	cd ./api_gateway && swag init -g **/**/*.go


# ==============================================================================
# Make local SSL Certificate

make_cert:
	echo "Generating SSL certificates"
	sh ./user/ssl/instructions.sh

# ==============================================================================
# Run linter for all services

linter:
	echo "Starting linters"
	cd comments && golangci-lint run ./...
	cd ..
	cd user && golangci-lint run ./...
	cd ..
	cd sessions && golangci-lint run ./...
	cd ..
	cd hotels && golangci-lint run ./...



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
# Postgresql migrations for microservices



# ==============================================================================
# Go migrate postgresql User service

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
# Go migrate postgresql Images service

images_dbname = images_db
images_port = 5434
images_SSL_MODE = disable

force_images_db:
	migrate -database postgres://postgres:postgres@localhost:$(images_port)/$(images_dbname)?sslmode=$(images_SSL_MODE) -path images/migrations force 1

version_images_db:
	migrate -database postgres://postgres:postgres@localhost:$(images_port)/$(images_dbname)?sslmode=$(images_SSL_MODE) -path images/migrations version

migrate_images_db_up:
	migrate -database postgres://postgres:postgres@localhost:$(images_port)/$(images_dbname)?sslmode=$(images_SSL_MODE) -path images/migrations up 1

migrate_images_db_down:
	migrate -database postgres://postgres:postgres@localhost:$(images_port)/$(images_dbname)?sslmode=$(images_SSL_MODE) -path images/migrations down 1


# ==============================================================================
# Go migrate postgresql Hotels service

hotels_dbname = hotels_db
hotels_port = 5435
hotels_SSL_MODE = disable

force_hotels_db:
	migrate -database postgres://postgres:postgres@localhost:$(hotels_port)/$(hotels_dbname)?sslmode=$(hotels_SSL_MODE) -path hotels/migrations force 1

version_hotels_db:
	migrate -database postgres://postgres:postgres@localhost:$(hotels_port)/$(hotels_dbname)?sslmode=$(hotels_SSL_MODE) -path hotels/migrations version

migrate_hotels_db_up:
	migrate -database postgres://postgres:postgres@localhost:$(hotels_port)/$(hotels_dbname)?sslmode=$(hotels_SSL_MODE) -path hotels/migrations up 1

migrate_hotels_db_down:
	migrate -database postgres://postgres:postgres@localhost:$(hotels_port)/$(hotels_dbname)?sslmode=$(hotels_SSL_MODE) -path hotels/migrations down 1



