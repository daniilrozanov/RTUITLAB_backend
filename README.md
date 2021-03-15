# RTUITLAB_backend

# Install and Run services

To install services, clone this repository in any folder, go into created folder by ```cd``` and run
```
  docker-compose pull
  docker-compose up
```
Servises run by pre-compiled binary file, so if you want to do any changes and up it, you need to execute ```CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd``` in each directory of ```d```
