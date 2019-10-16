DOCKER_TAG?=dev
ENVIRONMENT?=dev

######################
# Agent
######################
endpoint-runner-build:
	docker build -t managedkube/endpoint-runner:${DOCKER_TAG} --file ./endpoint-runner/Dockerfile .

endpoint-runner-push:
	docker push managedkube/endpoint-runner:${DOCKER_TAG}