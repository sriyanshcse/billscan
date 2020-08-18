NAME := okscan
DOCKER_REPO := gcr.io/okcredit-42/$(NAME)

###################
# TARGETS
#################
default: run

run:
	@go run cmd/main.go

build: OUT ?= $(NAME)
build:
	@CGO_ENABLED=0 go build -o $(OUT) cmd/main.go


docker:
	docker build --build-arg SSH_PRIVATE_KEY="$$(cat $(HOME)/.ssh/id_rsa)" -t $(NAME) -f Dockerfile ./..

docker-run: docker
docker-run:
	docker run -p "8080:8080" -v "$$(pwd)/config.yml:/config.yml" $(NAME)

docker-push: TAG ?= latest
docker-push: docker
	docker tag $(NAME) $(DOCKER_REPO):$(TAG)
	docker push $(DOCKER_REPO):$(TAG)