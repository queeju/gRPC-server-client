package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	transmitter "greaterm/alien_detector/gen/go"
	plog "greaterm/alien_detector/internal/slogpretty"
)

var (
	port   = flag.Int("port", 8888, "The server port")
	logger *slog.Logger
)

func main() {
	flag.Parse()

	logger = plog.SetupPrettySlog("stdout")

	logger.Info("Starting transmitter server", slog.Int("port", *port))

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logger.Error("Failed to listen: %v", err.Error(), "")
		return
	}
	grpcServer := grpc.NewServer()
	transmitterServer := &server{}
	transmitter.RegisterTransmitterServiceServer(grpcServer, transmitterServer)
	if err = grpcServer.Serve(lis); err != nil {
		logger.Error("Server error: %v", err.Error(), "")
	}
}
