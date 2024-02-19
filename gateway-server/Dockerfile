FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./gateway-server/*.go ./gateway-server/

COPY ./proto ./proto

WORKDIR /app/gateway-server
RUN go build -o /gateway-server

EXPOSE 8081

CMD ["/gateway-server"]
