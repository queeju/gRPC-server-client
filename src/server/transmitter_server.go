package main

import (
	"log/slog"
	"math/rand"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	transmitter "greaterm/alien_detector/gen/go"
)

const (
	minMean = -10
	maxMean = 10
	stdLow  = 0.3
	stdHigh = 1.5
	maxVals = 7500
)

type server struct {
	transmitter.UnimplementedTransmitterServiceServer
}

func (s *server) GetData(empty *empty.Empty, stream transmitter.TransmitterService_GetDataServer) error {
	var err error
	mean, sd := doMath()
	uuid := uuid.New().String()
	for i := 0; i < maxVals; i++ {
		freq := getFrequency(mean, sd)
		res := &transmitter.Response{
			SessionId: uuid,
			Frequency: freq,
			Time:      timestamppb.Now(),
		}
		if err = stream.Send(res); err != nil {
			logger.Info("Stream ended")
			break
		}
	}
	logger.Info("Session ID created", slog.String("ID", uuid))
	logger.Info("Normal distribution characteristics", slog.Float64("mean", mean), slog.Float64("SD", sd))
	return err
}

func doMath() (float64, float64) {
	mean := rand.Float64()*(maxMean-minMean) + minMean
	sd := rand.Float64()*(stdHigh-stdLow) + stdLow

	logger.Info("Normal distribution generated", slog.Float64("mean", mean), slog.Float64("SD", sd))
	return mean, sd
}

func getFrequency(mean, sd float64) float64 {
	freq := rand.NormFloat64()*sd + mean

	logger.Info("Frequency sampled", slog.Float64("frequency", freq))
	return freq
}
