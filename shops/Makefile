build_docker_image: ./main
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd
	docker build -t shops .
push_docker_image:
	docker tag shops ddzzan/shops
	docker push ddzzan/shops
all: build_docker_image push_docker_image