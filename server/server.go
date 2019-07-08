package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Ajod/Exercise/proto"

	"google.golang.org/grpc"
)

// These variables constitute the default command line modifiable connectivity parameters
var host = flag.String("host", "localhost", "The client host")
var port = flag.String("port", "8081", "The port the client should listen on")

type commandRunnerServer struct {
	storedNum []int64
}

func (s *commandRunnerServer) Compute(ctx context.Context, cp *pb.ComputeParam) (*pb.Result, error) {
	res := pb.Result{Command: "compute"}

	result := int64(0)
	if len(s.storedNum) > 1 {
		result = s.storedNum[0]
	}

	i := 0
	if result != 0 {
		i = 1
	}
	for ; i < len(s.storedNum); i++ {
		result -= s.storedNum[i]
	}
	s.storedNum = nil

	res.Results = append(res.Results, &pb.Arg{Name: "result", Value: result})
	return &res, nil
}

func (s *commandRunnerServer) Sub(ctx context.Context, val *pb.Value) (*pb.Result, error) {
	res := pb.Result{Command: "sub"}
	s.storedNum = append(s.storedNum, val.Value)
	res.Results = append(res.Results, &pb.Arg{Name: "pending_operations_nb", Value: int64(len(s.storedNum))})
	return &res, nil
}

func (s *commandRunnerServer) Sleep(ctx context.Context, val *pb.Value) (*pb.Result, error) {
	res := pb.Result{Command: "sleep"}

	t0 := time.Now()
	time.Sleep(time.Duration(val.Value) * time.Millisecond)
	t1 := time.Now()

	res.Results = append(res.Results, &pb.Arg{Name: "actual_duration", Value: int64(t1.Sub(t0) / time.Millisecond)})

	return &res, nil
}

func newServer() *commandRunnerServer {
	server := &commandRunnerServer{}
	return server
}

func main() {
	duck := flag.Bool("duck", false, "Shows a duck.")

	flag.Parse()

	if *duck == true {
		fmt.Printf("   _\n,_(')> Server Started.\n\\___)")
	}

	lis, err := net.Listen("tcp", *host+":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// We create the GRPC server
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	pb.RegisterCommandRunnerServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
