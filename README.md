# GO REST Service - Template

[![Go Report Card](https://goreportcard.com/badge/github.com/rakin92/go-service-template)](https://goreportcard.com/report/github.com/rakin92/go-service-template)

A simple Go REST service using:

- [Gin-gonic](https://gin-gonic.com) web framework
  - `go get -u github.com/gin-gonic/gin`
- [Goth](https://github.com/markbates/goth) for OAuth2 connections
  - `go get github.com/markbates/goth`
- [GORM](http://gorm.io) as DB ORM
  - `go get -u github.com/jinzhu/gorm`
  - [Gomigrate](https://gopkg.in/gormigrate.v2) orm migrations
    - `go get gopkg.in/gormigrate.v2`
- [Migrate](https://github.com/golang-migrate/migrate) for migrations scripts
  - `go get -u github.com/golang-migrate/migrate/v4`
- [Zerolog](https://github.com/rs/zerolog) for formatted logging
  - `go get -u github.com/rs/zerolog`

## Development with locally

Clone the `example.env` file to `.env` and update the values for your local development.
Run it locally with hot-reload:
```
sh scripts/run-air.sh
```

Run it locally without hot-reload:
```
sh scripts/run-dev.sh
```

## Development with docker

Just run it with `docker-compose`:

```
docker-compose run dev
```

And you'll have the service for your development and testing.

## Deployment

Use docker, swarm or kubernetes, GCP, AWS, DO, you name it.

Running `prod.dockerfile` will build a multistaged build that will give you a slim image containing just the service executable.

### With `docker-compose`

```
docker-compose build prod
```

or

```
docker-compose run prod
```

### Build from the `prod.dockerfile`

```
docker build -f docker/prod.dockerfile -t go-service.prod ./
```