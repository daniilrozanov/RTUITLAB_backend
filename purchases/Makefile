build_docker_image: ./main
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd
	docker build -t purchases .
push_docker_image:
	docker tag purchases ddzzan/purchases
	docker push ddzzan/purchases
all: build_docker_image push_docker_image