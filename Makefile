
export KO_DOCKER_REPO=docker.io/famargon

build:
	go build -o ./bin/app ./cmd/app.go

push:
	ko publish ./cmd
	ko publish ./cmd/timeline

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

deploy-redis:
	helm install redis bitnami/redis -n crud-app

.PHONY: deploy
deploy:
	kubectl create namespace crud-app | true
	$(MAKE) deploy-redis
	kubectl apply -f .dapr/configuration.yaml -n crud-app
	kubectl apply -f .dapr/components -n crud-app
	kubectl apply -f deploy -n crud-app

apply:
	kubectl apply -f .dapr/configuration.yaml -n crud-app
	kubectl apply -f .dapr/components -n crud-app
	kubectl apply -f deploy -n crud-app
