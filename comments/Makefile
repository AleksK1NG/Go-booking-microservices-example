# ==============================================================================
# Go migrate postgresql User service

comments_dbname =comments_db
comments_port = 5436
comments_SSL_MODE = disable

force_user_db:
	migrate -database postgres://postgres:postgres@localhost:$(comments_port)/$(comments_dbname)?sslmode=$(comments_SSL_MODE) -path migrations force 1

version_user_db:
	migrate -database postgres://postgres:postgres@localhost:$(comments_port)/$(comments_dbname)?sslmode=$(comments_SSL_MODE) -path migrations version

migrate_user_db_up:
	migrate -database postgres://postgres:postgres@localhost:$(comments_port)/$(comments_dbname)?sslmode=$(comments_SSL_MODE) -path migrations up 1

migrate_user_db_down:
	migrate -database postgres://postgres:postgres@localhost:$(comments_port)/$(comments_dbname)?sslmode=$(comments_SSL_MODE) -path migrations down 1

