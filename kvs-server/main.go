package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/btschwartz12/redis-lite/proto/kv_store"
	"google.golang.org/grpc"
)

type myKVStoreServer struct {
	kv_store.UnimplementedKeyValueStoreServer
}


func (s myKVStoreServer) Store(ctx context.Context, req *kv_store.StoreRequest) (*kv_store.StoreResponse, error) {
	log.Printf("Storing key: %s, value: %s", req.Key, req.Value)
	return &kv_store.StoreResponse {
		Success: true,
		Error: "none",
	}, nil
}

func (s myKVStoreServer) Retrieve(ctx context.Context, req *kv_store.RetrieveRequest) (*kv_store.RetrieveResponse, error) {
	log.Printf("Retrieving key: %s", req.Key)
	return &kv_store.RetrieveResponse {
		Value: "value",
		Error: "none",
	}, nil
}

func (s myKVStoreServer) Delete(ctx context.Context, req *kv_store.DeleteRequest) (*kv_store.DeleteResponse, error) {
	log.Printf("Deleting key: %s", req.Key)
	return &kv_store.DeleteResponse {
		Success: true,
		Error: "none",
	}, nil
}

func main() {
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = ":8080"
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	serverRegistrar := grpc.NewServer()
	service := &myKVStoreServer{}

	kv_store.RegisterKeyValueStoreServer(serverRegistrar, service)

	log.Println("Starting server on ", port)
	err = serverRegistrar.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}