package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"greaterm/alien_detector/analyse"
	transmitter "greaterm/alien_detector/gen/go"
	plog "greaterm/alien_detector/internal/slogpretty"
	"greaterm/alien_detector/pgdb"
)

const maxMessages = 5000

var (
	port  = flag.Int("port", 8888, "Server port")
	coef  = flag.Float64("k", 2.0, "Anomaly coefficient")
	setup = flag.Bool("s", false, "Table setup")
	pgConf = flag.String("pgconf", "postgres://postgres:123@localhost:5432/postgres?sslmode=disable", "Postgres address")
	logger *slog.Logger
)

func main() {
	flag.Parse()

	logger = plog.SetupPrettySlog("stdout")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", *port), opts...)
	if err != nil {
		logger.Error("Failed to dial:", err.Error(), "")
		return
	}
	defer conn.Close()

	client := transmitter.NewTransmitterServiceClient(conn)

	stream, err := client.GetData(context.Background(), nil)
	if err != nil {
		logger.Error("Client error:", err.Error(), "")
		return
	}

	logger.Info("Starting transmission", slog.Int("port", *port))

	a := analyse.NewAnalyser(maxMessages, logger, *coef)

	data, sum := a.ReadStream(stream)

	mean, sd := a.MakePrediction(data, sum)
	logger.Info("Predicted values", slog.Float64("mean", mean), slog.Float64("SD", sd))
	err = pgdb.DbConnect(*pgConf)
	if err != nil {
		logger.Error("Database connection error", err.Error(), "")
		return
	}

	err = pgdb.AddNewTable(logger, *setup)
	if err != nil {
		logger.Error("Database error", err.Error(), "")
		return
	}
	a.CatchAnomalies(mean, sd, stream)
}
