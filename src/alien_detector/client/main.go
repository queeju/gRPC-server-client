package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"sync"

	"gonum.org/v1/gonum/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	transmitter "greaterm/alien_detector/gen/go"
	// "greaterm/alien_detector/internal/mylog"
	plog "greaterm/alien_detector/internal/slogpretty"
)

const maxMessages = 200

var (
	port   = flag.Int("port", 8888, "Server port")
	k      = flag.Float64("k", 2.0, "Anomaly coefficient")
	logger *slog.Logger
)

func main() {
	flag.Parse()

	// TODO: handle err
	logger = plog.SetupPrettySlog("stdout")
	// logger, err := .SetupLogger("stdout")

	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Warn("warning")
	logger.Error("error msg!")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", *port), opts...)
	if err != nil {
		logger.Error("Failed to dial:", err)
	}
	defer conn.Close()

	client := transmitter.NewTransmitterServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := &transmitter.Request{Req: true}
	stream, err := client.Transmit(context.Background(), req)
	if err != nil {
		logger.Error("Client error:", err)
		return
	}

	logger.Info("Starting transmission", slog.Int("port", *port))

	// data := make(map[int32]float64, maxMessages)
	data := make([]float64, maxMessages)
	var (
		mu      sync.Mutex
		counter int
		sum     float64
	)

	go func() {
		for {
			response, err := stream.Recv()
			if err != nil {
				logger.Error("Streaming error:", err)
			}
			mu.Lock()
			// data[response.GetTime().Nanos] = response.GetFrequency()
			data[counter] = response.GetFrequency()
			counter++
			sum += response.GetFrequency()
			mu.Unlock()

			if counter >= maxMessages {
				cancel()
				return
			}
		}
	}()

	// Wait for the streaming to complete
	<-ctx.Done()
	total := len(data)
	logger.Info("Stream finished", slog.Int("messages", total))

	// mean := findMean(data)
	// sd := findSD(data)
	mean, sd := makePrediction(data, sum)
	logger.Info("Predicted values", slog.Float64("mean", mean), slog.Float64("SD", sd))

	anomalies := findAnomalies(data, mean, sd, *k)
	if anomalies > 0 {
		logger.Info("Anomalies found", slog.Int("total", anomalies))
	} else {
		logger.Info("No anomalies found")
	}
}

func findAnomalies(data []float64, mean, sd, k float64) int {
	var total int
	for _, val := range data {
		if math.Abs(val-mean) > sd*k {
			logger.Info("Found anomaly", slog.Float64("value", val))
			total++
		}
	}
	return total
}

func makePrediction(data []float64, sum float64) (float64, float64) {
	mean := sum / maxMessages
	sd := stat.StdDev(data, nil)
	return mean, sd
}

func findMean(data map[int32]float64) float64 {
	var res float64
	for _, val := range data {
		res += val
	}
	res /= float64(len(data))
	return res
}

func findSD(data map[int32]float64) float64 {
	arr := make([]float64, 100)
	i := 0
	for _, val := range data {
		arr[i] = val
		i++
	}
	return stat.StdDev(arr, nil)
}
