### Golang [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) microservices example with Prometheus, Grafana monitoring and Jaeger opentracing ⚡️

#### 👨‍💻 Full list what has been used:
* [GRPC](https://grpc.io/) - gRPC
* [RabbitMQ](https://github.com/streadway/amqp) - RabbitMQ
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [Docker](https://www.docker.com/) - Docker
* [Prometheus](https://prometheus.io/) - Prometheus
* [Grafana](https://grafana.com/) - Grafana
* [Jaeger](https://www.jaegertracing.io/) - Jaeger tracing
* [Go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware) - interceptor chaining, auth, logging, retries and more
* [Opentracing-go](https://github.com/opentracing/opentracing-go) - OpenTracing API for Go
* [Prometheus-go-client](https://github.com/prometheus/client_golang) - Prometheus instrumentation library for Go applications

#### Local development usage:
    make local // run all containers

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000

### RabbitMQ UI:

http://localhost:15672

### Swagger UI by default:

* https://localhost:8081/swagger/index.html - auth
* https://localhost:8016/swagger/index.html - gateway