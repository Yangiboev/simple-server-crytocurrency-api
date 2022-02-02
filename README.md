### Golang simple REST API🚀

## Initial thoughts

At the very beginning, I wanted to create as simple server. But, as the repo shows my ability to code, I decided to come up with better server.

* [echo](https://github.com/labstack/echo) - Web framework,  I have high desire to learn more about learning echo web framework, so decided to make it with echo. Honestly, it took a little time to adopt the framework but it is fun. 
* [viper](https://github.com/spf13/viper) - Go environment configuration 
* [go-redis](https://github.com/go-redis/redis) - Type-safe Redis client for Golang
* [zap](https://github.com/uber-go/zap) - Logger for logging
* [swag](https://github.com/swaggo/swag) - Swagger for documentation
* [Docker](https://www.docker.com/) - Docker for more simplicity in development
* [error-group](golang.org/x/sync/errgroup) - To reduce latency, concurrency

#### Recomendation for local development most comfortable usage:
    make local // run all containers
    
#### Recomendation for local development usage:    
    make run // it's easier way to attach debugger or rebuild/rerun project

#### 🙌👨‍💻🚀 Docker-compose files:
    docker-compose.yml - run redis

### Docker development usage:
    make docker

### Local development usage:
    make local
    make run

### SWAGGER UI:
https://localhost:8000/swagger/index.html
