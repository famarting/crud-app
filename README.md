# TODO list demo app

## Setup

```
skind start
dapr init -k
kubectl create namespace crud-app
kubectl apply -f .dapr/configuration.yaml
kubectl apply -f .dapr/components
kubectl apply -f deploy
```

https://docs.dapr.io/getting-started/configure-state-pubsub/