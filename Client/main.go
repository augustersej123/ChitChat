package client

import (
	"log"

	pb "github.com/augustersej123/ChitChat/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[Client] failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewChitChatClient(conn)
	_ = client
	log.Println("[Client] [Connect] connected to server")
}
