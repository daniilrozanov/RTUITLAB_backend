FROM golang:1.16-buster

ENV GOPATH=/
COPY ./ ./
RUN go mod download && \
go build -o main ./cmd/main.go

CMD ["./main"]