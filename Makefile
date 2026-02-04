.PHONY: install tests coverage docker-run docker-stop docker-restart docker-debug docker-clear docker-logs podman-run podman-stop podman-restart podman-debug podman-clear podman-logs

install:
	go mod tidy && go mod vendor

tests:
	go test -covermode=set ./... -coverprofile=coverage.txt && go tool cover -func=coverage.txt
coverage:
	go test -v -covermode=set ./... -coverprofile=coverage.txt && go tool cover -html=coverage.txt -o coverage.html && xdg-open coverage.html

docker-run:
	docker run --rm -d -p 80:80 -p 443:443 -v ./config.yaml:/app/config.yaml --add-host host.docker.internal:host-gateway --name proxier gouef/proxier

docker-stop:
	docker stop proxier

docker-restart:
	-@$(MAKE) docker-stop
	@$(MAKE) docker-run

docker-build:
	docker build -t gouef/proxier .

docker-debug:
	docker run -it gouef/proxier:latest /bin/sh

docker-clear:
	docker container rm proxier

docker-logs:
	docker logs proxier

#Podman
podman-run:
	sudo podman run --rm -d -v ./config.yaml:/app/config.yaml:Z --network=host --name proxier gouef/proxier

podman-stop:
	sudo podman stop proxier

podman-restart:
	-@$(MAKE) podman-stop
	@$(MAKE) podman-run

podman-build:
	podman build -t gouef/proxier .

podman-debug:
	sudo podman run -it gouef/proxier:latest /bin/sh

podman-clear:
	podman container rm proxier

podman-logs:
	sudo podman logs proxier