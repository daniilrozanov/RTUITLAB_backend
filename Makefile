all: rebuild_all_docker repush_all_docker

repush_all_docker:
	cd shops && make push_docker_image
	cd purchases && make push_docker_image
	cd fabric && make push_docker_image

rebuild_all_docker:
	cd shops && make build_docker_image
	cd purchases && make build_docker_image
	cd fabric && make build_docker_image

