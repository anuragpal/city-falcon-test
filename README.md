# City Falcon Test

## Requirement
[Initial Requirement](https://github.com/anuragpal/city-falcon-test/blob/main/docs/requirement.md)


## Setup Requirement
* System should have docker installed
* [Enable Postgres Stats](https://github.com/anuragpal/city-falcon-test/blob/main/docs/slow-query.md)

## Installation
### For Development

```sh
git clone https://github.com/anuragpal/city-falcon-test.git
cd city-falcon-test
docker-compose -f build-dev.yml up --build
```

** Verify Installation Working Property **
```sh
curl --location '127.0.0.1:3000'
```

** Running Test Cases ** 
Use new terminal
```sh
docker ps
```sh

** docker ps ** will provide you container id for `backend-api` use this id for loggin to system

```sh
docker exec -it 020723f2963b bash
```
Change 020723f2963b with docker id provided by system, Once logged in run command
```sh
go get github.com/stretchr/testify/assert
watch go test -v
```

** Running Test Coverage **
Use new terminal
```sh
docker ps
```sh

** docker ps ** will provide you container id for `backend-api` use this id for loggin to system

```sh
docker exec -it 020723f2963b bash
```
Change 020723f2963b with docker id provided by system, Once logged in run command
```sh
go get github.com/stretchr/testify/assert
watch go test -cover
```