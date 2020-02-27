.PHONY: all
GOOS=ddd
all: go docker kind-load helm

go:
	GOOS=linux CGO_ENABLED=0 go build -o nse ./cmd/nsed/main.go

docker:
	docker build -t cisco/cnns-nse:latest .

kind-load:
	kind load docker-image cisco/cnns-nse:latest

helm:
	helm install --name cnns-nse ./k8s/cnns-nse