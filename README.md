# Giftano Coding Test: Simple CRUD API with Golang

This is a simple CRUD API with golang for giftano coding test.
For the use case this app is using simple product management with categories.
The product category is a hierachical structure (tree) and a product can belong to 
multiple categories. 

This app is deployed with docker swarm in an AWS EC2 instance.

For API docs and functionalitty test, use [postman](https://documenter.getpostman.com/view/8923045/SWEB1FD4) 
(require desktop app installation).
To run functionality test, go to folder testing and sequentially execute the requests.

API URL: [http://ec2-54-169-218-62.ap-southeast-1.compute.amazonaws.com]

## Building the app in docker
```
docker build --build-arg SVC_NAME=product  --tag giftano-crud-golang/app-api:latest -f ./docker/Dockerfile .
```

## Run app in docker swarm
** You need to run this command on a swarm manager instance 
```
docker stack deploy -c docker/docker-compose.yml giftano-crud-golang
```

## Build and run locally
** Make sure you have mysql db instance ready and configure the environment variables (you can use env-dev.sh) **
```
source env-dev.sh
go run cmd/product.go
```