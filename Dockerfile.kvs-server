FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./kvs-server/*.go ./kvs-server/

COPY ./proto ./proto

WORKDIR /app/kvs-server
RUN go build -o /kvs-server

EXPOSE 8080

CMD ["/kvs-server"]
