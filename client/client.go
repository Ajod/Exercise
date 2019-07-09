package main

import (
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"flag"
	"fmt"
	"io/ioutil"
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

// These variables constitute the default command line modifiable connectivity parameters
var host = flag.String("host", "localhost", "The client host")
var backhost = flag.String("backhost", "localhost", "The backend host")
var port = flag.String("port", "4242", "The port the client should listen on")
var backport = flag.String("backport", "8081", "The port the client should connect to on the backend")

var (
	// Define the backend server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8081", "gRPC server endpoint")
)

type commandRequest struct {
	Command string `json:"command"`
	Args    []pb.Arg
}

func executes(writer http.ResponseWriter, req *http.Request) {
	// Get the json format of the request body
	cr := commandRequest{}
	bytestream, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	json.Unmarshal(bytestream, &cr)

	// Initiate a connection with the backend server
	conn, err := grpc.Dial(*backhost+":"+*backport, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// Declare a grpc client to make the calls, giving it the backend connection
	client = pb.NewCommandRunnerClient(conn)

	var res *pb.Result

	// Select the method corresponding to the command
	switch cr.Command {
	case "sub":
		res, err = client.Sub(context.Background(), &pb.Value{Value: cr.Args[0].Value})
	case "compute":
		res, err = client.Compute(context.Background(), &pb.ComputeParam{})
	case "sleep":
		res, err = client.Sleep(context.Background(), &pb.Value{Value: cr.Args[0].Value})
	default:
		fmt.Println("Unrecognized command: ", cr.Command)
	}
	if err != nil {
		panic(err)
	}

	// Put the command information back in result
	if cr.Args != nil {
		for _, arg := range cr.Args {
			res.Args = append(res.Args, &arg)
		}
	}
	// Build the json format of the response
	b, err := json.Marshal(res)

	writer.Write(b)
}

func run() error {
	// Start GRPC server
	return http.ListenAndServe(*host+":"+*port, router)
}

func main() {
	duck := flag.Bool("duck", false, "Shows a duck.")
	flag.Parse()

	if *duck == true {
		fmt.Printf("   _\n,_(')> API client Started.\n\\___)")
	}

	// Define available routes
	router.HandleFunc("/executes", executes).
		Methods("POST").
		Headers("Content-Type", "application/json").
		Headers("Accept", "application/json")

	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
