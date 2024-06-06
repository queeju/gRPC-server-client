package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	// "os"

	"google.golang.org/grpc"

	transmitter "greaterm/alien_detector/gen/go"
	// "greaterm/alien_detector/internal/mylog"
	plog "greaterm/alien_detector/internal/slogpretty"
)

var (
	port    = flag.Int("port", 8888, "The server port")
	logFile = flag.String("l", "stdout", "path/to/log")
	logger  *slog.Logger
)

func main() {
	flag.Parse()

	// TODO: handle err
	// logger, err := mylog.SetupLogger(*logFile)
	// logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger = plog.SetupPrettySlog(*logFile)

	logger.Info("Starting transmitter server", slog.Int("port", *port))

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logger.Error("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	transmitterServer := &server{}
	transmitter.RegisterTransmitterServiceServer(grpcServer, transmitterServer)
	if err = grpcServer.Serve(lis); err != nil {
		logger.Error("Server error: %v", err)
	}
}
