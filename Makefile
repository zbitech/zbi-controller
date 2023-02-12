#IMAGE_REPO=localhost:5001
#IMAGE_REPO=jakinyele
#VERSION=

compile:
	go build -v ./...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zbi-controller

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger swagger.yaml

unit_tests:
	go clean -testcache
#	export KUBECONFIG=$(HOME)/.kube/config &&
#	kind create cluster --name klient --kubeconfig=kubeconfig
	export ASSET_PATH_DIRECTORY=$(PWD)/../charts/controller/ && \
	export DATA_PATH=$(PWD)/./fake-zbi/data/ && \
	export ZBI_CONFIG_DIRECTORY=$(PWD)/../charts/controller/zbi-conf/ && \
	export ZBI_TEMPLATE_DIRECTORY=$(PWD)/../charts/controller/zbi-templates/ && \
	export KUBERNETES_IN_CLUSTER=false && \
	export KUBECONFIG=$(PWD)/kubeconfig && \
	go test -v ./...
#	kind delete cluster --name klient

run:
#	export ASSET_PATH_DIRECTORY=$(PWD)/../control-plane/zbi/conf
	export ASSET_PATH_DIRECTORY=$(PWD)/charts/controller/ && \
	export ZBI_CONFIG_DIRECTORY=$(PWD)/charts/controller/zbi-conf/ && \
	export ZBI_TEMPLATE_DIRECTORY=$(PWD)/charts/controller/zbi-templates/ && \
	export KUBERNETES_IN_CLUSTER=false && \
	export KUBECONFIG=$(HOME)/.kube/config && \
	go run main.go --port 8180

image:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o zbi-controller
#	export VERSION=$(git rev-parse HEAD)
	docker build -t ${ZBI_REPO}/zbi-controller:latest -t ${ZBI_REPO}/zbi-controller:${VERSION} .


push_ecr:
	aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${ZBI_REPO}
#	export VERSION=$(git rev-parse HEAD)
	docker push ${ZBI_REPO}/zbi-controller:${VERSION}
	docker push ${ZBI_REPO}/zbi-controller:latest
