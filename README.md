# Redis-Lite

This is an extremely simple implementation of a key-value store server written in Go that communicates over gRPC, with a JSON HTTP frontend.

Since it is somewhat similar to Redis, I decided to call it Redis-Lite, but this was just made for learning purposes.

## Quickstart

To start the services, you can use Docker Compose:
```bash
$ docker-compose up
```
The gateway server (with the JSON API) will be available at `http://localhost:8081`.

## Services

#### 1. Key-Value Store Server

This is the main service that stores the key-value pairs and communicates over gRPC. You can find the protobuf definition in [kv_store.proto](./proto/kv_store.proto), and the source code in the [kvs-server](./kv_store) directory.

At the moment, the server only supports the following operations:
- `Store` - Store a key-value pair
- `Retrieve` - Retrieve a value given a key
- `Delete` - Delete a key-value pair

I could have used an in-memory map to store the key-value pairs, but just for fun I used an actual Redis instance as the backend database! Although this completely defeats the purpose of this being a "lite" version of Redis...yeah I don't really have an excuse for this.

This component is containerized inside its [Dockerfile](./kv_store/Dockerfile).

#### 2. Gateway Server

The frontend of this "Redis-Lite" shouldn't have to communicate over gRPC, so I made a simple JSON API that acts as a 'gateway' to the KVS server. You can find the source code in the [gateway-server](./gateway-server) directory.

The frontend supports the same operations as the KVS server, but over HTTP:
- `POST /store` - Store a key-value pair
- `GET /retrieve` - Retrieve a value given a key
- `DELETE /delete` - Delete a key-value pair

This component is containerized inside its [Dockerfile](./gateway-server/Dockerfile), and it is the only service that is exposed to the host machine.

## Testing

To manually test these services, you simply need to start them locally, and use a tool like `curl` or Postman to send requests to the gateway server.

To start all the services, use Docker Compose:
```bash
$ docker-compose up
```

Then, you can use `curl` to send requests to the gateway server:
```bash
$ curl  --location 'http://localhost:8081/store' \
        --header 'Content-Type: application/json' \
        --data '{"key":"message", "value":"i love dogs"}'
```

```bash
$ curl  --location --request GET 'http://localhost:8081/retrieve' \
        --header 'Content-Type: application/json' \
        --data '{"key":"message"}'
```

```bash
$ curl  --location --request DELETE 'http://localhost:8081/delete' \
        --header 'Content-Type: application/json' \
        --data '{"key":"message"}'
```


I also created a Postman collection that you can use to test the services, although it is not too different than running the curl commands. You can find it here:

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://app.getpostman.com/run-collection/11432210-a5faef47-6542-4084-b5bb-ecea994375e2?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D11432210-a5faef47-6542-4084-b5bb-ecea994375e2%26entityType%3Dcollection%26workspaceId%3D7702f1b1-c090-48dc-8181-53754335f667)

## Further Improvements

There are many things that could be improved in this project, such as:

- Add more operations to the KVS server, such as `List` and `Update`
- Add a proper backend for the database
- Add an authentication scheme, such as JWT
- Add automated tests
