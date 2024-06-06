package main

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"
	rnd "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	transmitter "greaterm/alien_detector/gen/go"
)

var normalDist distuv.Normal

type server struct {
	transmitter.UnimplementedTransmitterServiceServer
}

func (s *server) Transmit(
	req *transmitter.Request,
	stream transmitter.TransmitterService_TransmitServer) error {
	doMath() // generate normal distribution
	uuid := uuid.New().String()
	logger.Info("Session ID created", slog.String("ID", uuid))
	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			time.Sleep(time.Second/10)
			freq := getFrequency()
			res := &transmitter.Response{
				SessionId: uuid,
				Frequency: freq,
				Time:      timestamppb.Now(),
			}
			if err := stream.SendMsg(res); err != nil {
				return err
			}
		}
	}
}

func doMath() {
	// mean from [-10, 10] interval
	mean := rand.Float64()*20 - 10
	// standard deviation from [0.3, 1.5].
	sd := rand.Float64()*1.2 + 0.3

	// Create a normal distribution with the specified mean and standard deviation
	normalDist = distuv.Normal{
		Mu:    mean,
		Sigma: sd,
		Src:   rnd.NewSource(uint64(time.Now().UnixNano())),
	}
	logger.Info("Normal distribution generated", slog.Float64("mean", mean), slog.Float64("SD", sd))
}

func getFrequency() float64 {
	// Generate a sample from the normal distribution
	freq := normalDist.Rand()

	logger.Info("Frequency sampled", slog.Float64("frequency", freq))
	return freq
}
