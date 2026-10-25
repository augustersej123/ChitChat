package main

import (
	"fmt"
	"io"
	"log"
	"unicode/utf8"

	pb "github.com/augustersej123/ChitChat/grpc"
)

type server struct {
	pb.UnimplementedChitChatServer
	broadcaster *Broadcaster
}

func (s *server) Chat(stream pb.ChitChat_ChatServer) error {
	//The first event must be a join
	first, err := stream.Recv()
	if err != nil {
		return err
	}

	join := first.GetJoin()
	if join == nil {
		return fmt.Errorf("First event must be a join")
	}
	// Reject empty participant names
	if join.Participant == "" {
		return fmt.Errorf("participant name must not be empty")
	}

	p := newParticipant(join.Participant)

	//Sender goroutine: forwards broadcasts to this client
	//It stops when the broadcaster closes p.Out on leave

	go func() {
		for msg := range p.Out { //receive messages the broadcaster delivers to this client
			if err := stream.Send(msg); err != nil { //If an error occurs
				log.Printf("[Server] event=send-error client=%s err=%v", p.ID, err)
			}
		}
	}() //() calls the function. go runs it as a goroutine

	log.Printf("[Server] event=client-connect client=%s", p.ID)
	defer log.Printf("[Server] event=client-disconnect client=%s", p.ID)

	s.broadcaster.Send(Event{Type: EventJoin, Participant: p, ClientTime: first.Timestamp})

	//Receive loop: handles the events from this client
	for {
		ev, err := stream.Recv()
		if err == io.EOF || err != nil {
			//Connection dropped without a leave. Local server event
			s.broadcaster.Send(Event{Type: EventLeave, Participant: p})
			return nil
		}

		switch e := ev.Event.(type) {
		case *pb.ClientEvent_Publish:

			if !utf8.ValidString(e.Publish.Content) {
				log.Printf("[Server] event=reject-publish client=%s reason=invalid-utf8", p.ID)
				continue
			}
			if utf8.RuneCountInString(e.Publish.Content) > 128 {
				log.Printf("[Server] event=reject-publish client=%s reason=too-long", p.ID)
				continue
			}
			s.broadcaster.Send(Event{Type: EventPublish, Participant: p, Content: e.Publish.Content, ClientTime: ev.Timestamp})

		case *pb.ClientEvent_Leave:
			s.broadcaster.Send(Event{Type: EventLeave, Participant: p, ClientTime: ev.Timestamp})
			return nil
		}
	}
}
