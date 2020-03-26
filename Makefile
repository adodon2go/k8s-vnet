GO_TOP := $(shell echo ${GOPATH} | cut -d ':' -f1)
export GO_TOP
export OUT_DIR=$(GO_TOP)/out
DOCKER_BUILD_TOP:=$(GO_TOP)/out/linux_amd64/debug/docker_build
.PHONY: all

all: go docker

go:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -o ${GO_TOP}/bin/linux_amd64/vl3_nse ./cmd/nsed/...

docker: go
	mkdir -p ${DOCKER_BUILD_TOP}/vl3/etc
	cp -p ./Dockerfile ${DOCKER_BUILD_TOP}/vl3/
	cp -p ${GOPATH}/bin/linux_amd64/vl3_nse ${DOCKER_BUILD_TOP}/vl3/
	cp -pr ./static/etc/* ${DOCKER_BUILD_TOP}/vl3/etc/
	cd ${DOCKER_BUILD_TOP}/vl3 && docker build -t ${ORG}/vl3-nse:${TAG} -f Dockerfile .
