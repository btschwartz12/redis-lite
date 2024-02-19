package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/btschwartz12/redis-lite/proto/kv_store"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type myKVStoreServer struct {
	kv_store.UnimplementedKeyValueStoreServer
	redisClient *redis.Client
}


func logRedisErrorAndGetFriendlyMessage(ctx context.Context, key string, err error) (string) {
	log.Printf("Error with key: %s, error: %s", key, err.Error())
	if err == redis.Nil {
		return "key " + key + " not found"
	}
	return "an unknown error occurred"
}


func (s myKVStoreServer) Store(ctx context.Context, req *kv_store.StoreRequest) (*kv_store.StoreResponse, error) {
	log.Printf("Storing value for key: %s", req.Key)
	err := s.redisClient.Set(ctx, req.Key, req.Value, 0).Err()
	if err != nil {
		errorMsg := logRedisErrorAndGetFriendlyMessage(ctx, req.Key, err)
		return &kv_store.StoreResponse {
			Success: false,
			Error: errorMsg,
		}, nil
	}
	log.Printf("Stored value for key: %s", req.Key)
	return &kv_store.StoreResponse {
		Success: true,
		Error: "none",
	}, nil
}

func (s myKVStoreServer) Retrieve(ctx context.Context, req *kv_store.RetrieveRequest) (*kv_store.RetrieveResponse, error) {
	log.Printf("Retrieving value for key: %s", req.Key)
	val, err := s.redisClient.Get(ctx, req.Key).Result()
	if err != nil {
		errorMsg := logRedisErrorAndGetFriendlyMessage(ctx, req.Key, err)
		return &kv_store.RetrieveResponse {
			Found: false,
			Value: "",
			Error: errorMsg,
		}, nil
	}
	log.Printf("Retrieved value for key: %s", req.Key)
	return &kv_store.RetrieveResponse {
		Found: true,
		Value: val,
		Error: "none",
	}, nil
}

func (s myKVStoreServer) Delete(ctx context.Context, req *kv_store.DeleteRequest) (*kv_store.DeleteResponse, error) {
	log.Printf("Deleting key: %s", req.Key)
	deleted, err := s.redisClient.Del(ctx, req.Key).Result()
	if err != nil {
		errorMsg := logRedisErrorAndGetFriendlyMessage(ctx, req.Key, err)
		return &kv_store.DeleteResponse {
			Success: false,
			Error: errorMsg,
		}, nil
	}
	if deleted == 0 {
        errorMsg := "Key " + req.Key + " not found"
        return &kv_store.DeleteResponse{
            Success: false,
            Error:   errorMsg,
        }, nil
    }
	log.Printf("Deleted key: %s", req.Key)
    return &kv_store.DeleteResponse{
        Success: true,
        Error:   "none",
    }, nil
}

func NewRedisClient() *redis.Client {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		redisAddress = "localhost:6379"
	}
    rdb := redis.NewClient(&redis.Options{
        Addr:     redisAddress,
        Password: "",
        DB:       0,
    })

    return rdb
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
	service := &myKVStoreServer{
		redisClient: NewRedisClient(),
	}

	kv_store.RegisterKeyValueStoreServer(serverRegistrar, service)

	log.Println("Starting server on ", port)
	err = serverRegistrar.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}