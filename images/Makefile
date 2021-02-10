# ==============================================================================
# Go migrate postgresql Images service

images_dbname = images_db
images_port = 5434
images_SSL_MODE = disable

force_images_db:
	migrate -database postgres://postgres:postgres@localhost:$(images_port)/$(images_dbname)?sslmode=$(images_SSL_MODE) -path migrations force 1

version_images_db:
	migrate -database postgres://postgres:postgres@localhost:$(images_port)/$(images_dbname)?sslmode=$(images_SSL_MODE) -path migrations version

migrate_images_db_up:
	migrate -database postgres://postgres:postgres@localhost:$(images_port)/$(images_dbname)?sslmode=$(images_SSL_MODE) -path migrations up 1

migrate_images_db_down:
	migrate -database postgres://postgres:postgres@localhost:$(images_port)/$(images_dbname)?sslmode=$(images_SSL_MODE) -path migrations down 1
