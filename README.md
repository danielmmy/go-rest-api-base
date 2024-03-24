# Go REST-api base

Skeleton project for go rest api

## Setup

A makefile is available at the root. To install [make](https://www.gnu.org/software/make/manual/make.html) on windows use [Chocolatey](https://chocolatey.org/) and running as admin execute:
```bash
choco install make
```

To run the api in docker you will need [Docker](https://www.docker.com/get-started/)

There are sample postman request under "postman_test_collection"  [folder](https://github.com/danielmmy/go-rest-api-base/tree/main/postman_test_collection)


## Use

### Building and start
```bash
    make up
```
### Stop
```bash
    make down
```
### Run locally
```bash
    make run
```
### Run test
```bash
    make test
```
### Check Coverage
```bash
    make cover
```

### Additional notes

The local server listen on port :8080 while the docker server listen on :15006
