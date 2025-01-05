package main

import (
	"net"
	pb "server/gen/proto/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type userService struct {
	pb.UnimplementedUserServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":5001")

	if err != nil {
		panic(err)
	}

	//NewServer creates a gRPC server which has no service registered
	// creating an instance of the server
	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, &userService{})

	//exposing gRPC service to be tested by postman
	reflection.Register(s)

	err = s.Serve(listener) // run the gRPC server
	if err != nil {
		panic(err)
	}
}
