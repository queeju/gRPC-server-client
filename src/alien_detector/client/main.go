package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	transmitter "greaterm/alien_detector/gen/go"
)

const maxMessages = 100

var (
	port = flag.Int("port", 8888, "The server port")
)

func main() {
	flag.Parse()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", *port), opts...)
	if err != nil {
		log.Fatalln("Failed to dial:", err)
	}
	defer conn.Close()

	client := transmitter.NewTransmitterServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &transmitter.Request{Req: true}
	stream, err := client.Transmit(context.Background(), req)
	if err != nil {
		log.Fatalln("Client error:", err)
	}

	data := make(map[int32]float64, maxMessages)
	var mu sync.Mutex
	counter := 0

	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				log.Fatalln("Streaming error:", err)
			}
			mu.Lock()
			data[response.GetTime().Nanos] = response.GetFrequency()
			counter++
			mu.Unlock()

			if counter >= maxMessages {
				cancel()
				return
			}
		}
	}()

	// Wait for the streaming to complete
	<-ctx.Done()
}
