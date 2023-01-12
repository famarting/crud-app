
export KO_DOCKER_REPO=docker.io/famargon

build:
	go build -o ./bin/app ./cmd/app.go

push:
	ko publish ./cmd
	ko publish ./cmd/timeline
	ko publish ./cmd/datagen
	ko publish ./cmd/errorgen
	ko publish ./cmd/consumer
	ko publish ./cmd/publisher
	ko publish ./cmd/service-a
	ko publish ./cmd/service-b
	ko publish ./cmd/service-c

rollout:
	kubectl delete pod -l app=crud-app
	kubectl delete pod -l app=timeline-app

crud-app-logs:
	k logs -l app=crud-app -c crud-app

run-mongo:
#-e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin
	docker run -it -p 27017:27017 mongo:4.0-xenial

dapr-run:
	dapr run --app-id crud-app --app-port 8080 --dapr-http-port 3500 ./bin/app serve -connStr dapr

setup-zipkin: deploy-zipkin
	skind expose zipkin

deploy-zipkin:
	kubectl create deployment zipkin --image openzipkin/zipkin
	kubectl expose deployment zipkin --type ClusterIP --port 9411

APP_NAMESPACE ?= crud-app

deploy-redis:
	helm upgrade --install redis bitnami/redis -n ${APP_NAMESPACE} --set architecture=standalone

deploy-redis-with-replication:
	helm upgrade --install redis bitnami/redis -n ${APP_NAMESPACE} --set replica.replicaCount=1

.PHONY: deploy
deploy:
	kubectl create namespace ${APP_NAMESPACE} | true
	$(MAKE) deploy-redis
	kubectl apply -f .dapr/configuration.yaml -n ${APP_NAMESPACE}
	kubectl apply -f .dapr/components -n ${APP_NAMESPACE}
	kubectl apply -f deploy -n ${APP_NAMESPACE}

apply:
	kubectl apply -f .dapr/configuration.yaml -n crud-app
	kubectl apply -f .dapr/components -n crud-app
	kubectl apply -f deploy -n crud-app
