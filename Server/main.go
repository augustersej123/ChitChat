package main

import (
	"log"
	"net"

	pb "github.com/augustersej123/ChitChat/grpc"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChitChatServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("[Server] failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChitChatServer(s, &server{})
	log.Println("[Server] [Startup] listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[Server] failed to serve: %v", err)
	}
}
