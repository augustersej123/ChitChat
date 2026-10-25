package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/augustersej123/ChitChat/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//read the clients name with a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("usage: go run ./client <name>")
		os.Exit(1)
	}
	name := os.Args[1] //os.Args[0] is the program itself, [1] is the first argument

	// log handler goes into main.
	f, err := os.OpenFile("client_"+name+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	// os.O_CREATE If file doesnt exist we create it
	//os.O_Wronly Opens file for write only
	//os.O_Append makes every write go to the end of the document, so new content doesnt overwrite from the top
	if err != nil {
		log.Fatalf("could not open the log file:%v", err)
		os.Exit(1)
	}

	defer f.Close()

	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds) //Timestamp format to microseconds
	log.SetPrefix("[Client " + name + "] ")

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(" failed to connect: %v", err)
		fmt.Printf("could not connect: %v\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewChitChatClient(conn)
	log.Printf("event=connect client=%s", name)

	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf(" failed to open the chat stream: %v", err)
		fmt.Printf("could not open the chat stream: %v\n", err)
		return
	}

	clock := &LogicalClock{}
	done := make(chan struct{})
	go receive(stream, clock, done) // Start receiving messages in a separate goroutine

	if err := sendJoin(stream, clock, name); err != nil {
		log.Fatalf("[Client] failed to join %v", err)
		fmt.Printf("failed to join: %v\n", err)
		return
	}

	fmt.Printf("Joined as %s. Type a message, or /leave to quit.\n", name)

	readInput(stream, clock)

	stream.CloseSend()
	<-done

	log.Printf("event=disconnect client=%s", name)
	fmt.Println("Left Chit Chat.")
}
