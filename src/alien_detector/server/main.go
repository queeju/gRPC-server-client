package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
	transmitter "greaterm/alien_detector/gen/go"
	plog "greaterm/alien_detector/internal/slogpretty"
)

var (
	port    = flag.Int("port", 8888, "The server port")
	logFile = flag.String("l", "", "path/to/log")
	logger  *slog.Logger
)

func main() {
	flag.Parse()
	logger = setupLogger()
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

func setupLogger() *slog.Logger {
	var lf *os.File
	if *logFile == "" {
		lf = os.Stdout
	} else {
		// TODO: implement logging into file
		lf = os.Stderr
	}

	opts := plog.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(lf)
	return slog.New(handler)
}
