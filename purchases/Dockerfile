FROM golang:1.16-alpine

ENV GOPATH=/
COPY ./ ./
#RUN apt-get update && apt-get -y install postgresql-client && \

RUN go mod download && \
go build -o main ./cmd/main.go

CMD ["./main"]


