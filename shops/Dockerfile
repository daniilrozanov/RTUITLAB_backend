#FROM golang:1.16-alpine

#ENV GOPATH=/
#COPY ./ ./
#RUN go mod download && \
#go build -o main ./cmd/main.go

FROM scratch

COPY ./configs ./configs
COPY main ./

CMD ["./main"]