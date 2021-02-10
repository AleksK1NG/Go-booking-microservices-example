# ==============================================================================
# Go migrate postgresql Hotels service

hotels_dbname = hotels_db
hotels_port = 5435
hotels_SSL_MODE = disable

force_hotels_db:
	migrate -database postgres://postgres:postgres@localhost:$(hotels_port)/$(hotels_dbname)?sslmode=$(hotels_SSL_MODE) -path migrations force 1

version_hotels_db:
	migrate -database postgres://postgres:postgres@localhost:$(hotels_port)/$(hotels_dbname)?sslmode=$(hotels_SSL_MODE) -path migrations version

migrate_hotels_db_up:
	migrate -database postgres://postgres:postgres@localhost:$(hotels_port)/$(hotels_dbname)?sslmode=$(hotels_SSL_MODE) -path migrations up 1

migrate_hotels_db_down:
	migrate -database postgres://postgres:postgres@localhost:$(hotels_port)/$(hotels_dbname)?sslmode=$(hotels_SSL_MODE) -path migrations down 1

