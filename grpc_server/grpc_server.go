package main

import (
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"flag"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	pb "github.com/Ajod/Exercise/proto" // Update
	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var conn *grpc.ClientConn
var client pb.CommandRunnerClient

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8081", "gRPC server endpoint")
)

func executes(writer http.ResponseWriter, req *http.Request) {
	//fmt.Println(req)
	fmt.Println("Sir, we have received a request on /executes.")
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = pb.NewCommandRunnerClient(conn)
	res, err := client.Sub(context.Background(), &pb.Command{Command: "sub", Args: nil})
	fmt.Println(err)
	bytestream, err := json.Marshal(res)
	writer.Write(bytestream)
	writer.WriteHeader(200)
	fmt.Println("Execute Processed")
}

func run() error {
	//ctx := context.Background()
	//ctx, cancel := context.WithCancel(ctx)
	//defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible

	//mux := runtime.NewServeMux()
	//	opts := []grpc.DialOption{grpc.WithInsecure()}
	//	err := gw.RegisterCommandRunnerHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	//	if err != nil {
	//		return err
	//	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":4242", router)
}

func main() {
	flag.Parse()

	router.HandleFunc("/executes", executes).
		Methods("POST").
		Headers("Content-Type", "application/json").
		Headers("Accept", "application/json")

	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
