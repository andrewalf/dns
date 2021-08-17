.PHONY: vendor
vendor:
	go mod vendor
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	docker pull golang:1.16 && docker pull alpine:3.14 && docker build -t dns_image -f ./build/Dockerfile .

.PHONY: run
run:
	docker run -d -e DNS_PORT=$(port) -e SECTOR_ID=$(sector_id) \
	--name dns_container_$(port) -p 127.0.0.1:$(port):$(port)/tcp dns_image

.PHONY: stop
stop:
	docker stop dns_container_$(port) && docker rm dns_container_$(port)

.PHONY: thanos_snapd
thanos_snap:
	docker stop $(shell docker ps -a -q --filter="name=dns_container_") && \
	docker rm $(shell docker ps -a -q --filter="name=dns_container_") && \
	docker rmi $(shell docker images -q dns_image) && \
	docker rmi $(shell docker images -q alpine:3.14) && \
	docker rmi $(shell docker images -q golang:1.16)
