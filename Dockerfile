FROM golang:latest

ENV GOROOT  /usr/local/go

WORKDIR /go/src/CRUDServer

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY /internal/configs/config.env config.env
COPY *.go ./


RUN go build -o /docker_crudserver

EXPOSE 8081

CMD ["/docker_crudserver"]