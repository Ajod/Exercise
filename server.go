package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/Ajod/Exercise/proto"

	"google.golang.org/grpc"
)

type commandRunnerServer struct {
	storedNum []int64
}

/*func (s *commandRunnerServer) Execute(ctx context.Context, cmd *pb.Command) (*pb.Result, error) {
	res := pb.Result{}
	if cmd.Command == "sub" {
		res.Command = "sub"
		for _, arg := range cmd.Args {
			s.storedNum = append(s.storedNum, arg.Value)
			res.Args = append(res.Args, arg)
			res.Results = append(res.Results, &pb.Arg{Name: "pending_operations_nb", Value: int64(len(s.storedNum))})
		}
	}
	fmt.Println("[OSCAR]: \"All is well Sir.\"")
	fmt.Println("   _\n,_(')>\n\\___)")
	return &pb.Result{}, nil
}*/

func (s *commandRunnerServer) Compute(ctx context.Context, cmd *pb.Command) (*pb.Result, error) {
	return nil, nil
}

func (s *commandRunnerServer) Sub(ctx context.Context, cmd *pb.Command) (*pb.Result, error) {
	fmt.Println("AAAAAAAAAAAAAAAAHHHHHHHH")
	return &pb.Result{}, errors.New("Subbed")
}

func (s *commandRunnerServer) Sleep(ctx context.Context, cmd *pb.Command) (*pb.Result, error) {
	return nil, nil
}

func newServer() *commandRunnerServer {
	server := &commandRunnerServer{}
	return server
}

func main() {
	fmt.Println("[OSCAR]: Starting the server Sir.")

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8081))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// We create the GRPC server
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	pb.RegisterCommandRunnerServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
