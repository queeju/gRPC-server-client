package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"

	transmitter "greaterm/alien_detector/gen/go"
)

var (
	port    = flag.Int("port", 8888, "The server port")
	logFile = flag.String("l", "", "path/to/log")
)

func main() {
	flag.Parse()
	logger := setupLogger()
	logger.Info("Starting transmitter server", slog.Int("port", *port))
	logger.Debug("debug msg")
	logger.Error("error msg")
	logger.Info("info msg")
	logger.Warn("warn msg")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logger.Error("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	transmitterServer := &server{}
	transmitter.RegisterTransmitterServiceServer(grpcServer, transmitterServer)
	grpcServer.Serve(lis)
}

func setupLogger() *slog.Logger {
	var lf *os.File
	if *logFile == "" {
		lf = os.Stdout
	} else {
		lf = os.Stderr
	}
	// TODO: add log file path handling
	return slog.New(slog.NewTextHandler(lf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
