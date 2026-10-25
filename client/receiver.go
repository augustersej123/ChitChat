package main

import (
	"fmt"
	"io"
	"log"

	pb "github.com/augustersej123/ChitChat/grpc"
)

// prints every message the server broadcasts and runs as a goroutine
func receive(stream pb.ChitChat_ChatClient, clock *LogicalClock, done chan<- struct{}) {
	defer close(done)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("event=stream-closed")
			return
		}
		if err != nil {
			log.Printf("event=receive-error err=%v", err)
			return
		}

		t := clock.Update(msg.Timestamp) // Lamport rule: max(local, received) + 1

		log.Printf("event=receive sender=%s content=%q msg_lamport=%d lamport=%d",
			msg.Sender, msg.Content, msg.Timestamp, t)
		fmt.Printf("[%d] %s: %s\n", msg.Timestamp, msg.Sender, msg.Content)
	}
}
