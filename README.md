# go-microservice-api

## RabbitMQ

```bash
docker build -t rabbitmq:latest -f docker/rabbitmq/Dockerfile .

kubectl create deployment rabbitmq --image=rabbitmq:latest

kubectl expose deployment rabbitmq --type=NodePort --port=5672
```

## API-microservice

```bash
docker build -t app-api:api -f docker/deploy/Dockerfile .

kubectl apply -f kubernetes/deployment.yaml

kubectl apply -f kubernetes/service.yaml
```
