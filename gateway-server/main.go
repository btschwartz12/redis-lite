package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/btschwartz12/redis-lite/proto/kv_store"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type gRPCAction func(client kv_store.KeyValueStoreClient, ctx context.Context, req proto.Message) (proto.Message, error)

func createGRPCAction(action gRPCAction, req proto.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverAddress := os.Getenv("GRPC_SERVER_ADDRESS")
		if serverAddress == "" {
			serverAddress = ":8080"
		}

		conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
		if err != nil {
			http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		client := kv_store.NewKeyValueStoreClient(conn)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		err = protojson.Unmarshal(body, req)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		log.Printf("Received request type: %T", req)

		resp, err := action(client, context.Background(), req)
		if err != nil {
			http.Error(w, "gRPC action failed", http.StatusInternalServerError)
			return
		}
		
		jsonData, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(resp)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

func storeHandler() http.HandlerFunc {
	f := func (client kv_store.KeyValueStoreClient, ctx context.Context, req proto.Message) (proto.Message, error) {
        return client.Store(ctx, req.(*kv_store.StoreRequest))
    }
	return createGRPCAction(f, &kv_store.StoreRequest{})
}

func retrieveHandler() http.HandlerFunc {
    f := func (client kv_store.KeyValueStoreClient, ctx context.Context, req proto.Message) (proto.Message, error) {
		return client.Retrieve(ctx, req.(*kv_store.RetrieveRequest))
	}
	return createGRPCAction(f, &kv_store.RetrieveRequest{})
}

func deleteHandler() http.HandlerFunc {
    f := func (client kv_store.KeyValueStoreClient, ctx context.Context, req proto.Message) (proto.Message, error) {
		return client.Delete(ctx, req.(*kv_store.DeleteRequest))
	}
	return createGRPCAction(f, &kv_store.DeleteRequest{})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/store", storeHandler()).Methods("POST")
    r.HandleFunc("/retrieve", retrieveHandler()).Methods("GET")
    r.HandleFunc("/delete", deleteHandler()).Methods("DELETE")

	log.Println("Starting server on :8081")
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = ":8081"
	}
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
