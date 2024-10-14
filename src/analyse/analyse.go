package analyse

import (
	"log/slog"
	"math"
	"time"

	"gonum.org/v1/gonum/stat"

	transmitter "greaterm/alien_detector/gen/go"
	"greaterm/alien_detector/pgdb"
)

type Analyser struct {
	maxMessages int
	logger      *slog.Logger
	coef        float64
}

func NewAnalyser(maxMess int, log *slog.Logger, coefficient float64) *Analyser {
	ob := new(Analyser)
	ob.maxMessages = maxMess
	ob.logger = log
	ob.coef = coefficient
	return ob
}

func (a Analyser) AnalyseStream(stream transmitter.TransmitterService_GetDataClient) {
	data, sum := a.ReadStream(stream)
	mean, sd := a.MakePrediction(data, sum)
	a.logger.Info("Predicted values", slog.Float64("mean", mean), slog.Float64("SD", sd))
	a.CatchAnomalies(mean, sd, stream)
}

func (a Analyser) ReadStream(stream transmitter.TransmitterService_GetDataClient) ([]float64, float64) {
	data := make([]float64, a.maxMessages)
	var (
		counter int
		sum     float64
	)
	for i := 0; i < a.maxMessages; i++ {
		response, err := stream.Recv()
		if err != nil {
			a.logger.Error("Streaming error:", err.Error(), "")
		}
		data[counter] = response.GetFrequency()
		counter++
		sum += response.GetFrequency()
	}
	return data, sum
}

func (a Analyser) MakePrediction(data []float64, sum float64) (float64, float64) {
	mean := sum / float64(a.maxMessages)
	sd := stat.StdDev(data, nil)
	return mean, sd
}

func (a Analyser) CatchAnomalies(mean, sd float64, stream transmitter.TransmitterService_GetDataClient) {
	for {
		time.Sleep(time.Second / 100)
		response, err := stream.Recv()
		if err != nil {
			a.logger.Error("Streaming error:", err.Error(), "")
		}
		freq := response.GetFrequency()
		if math.Abs(freq-mean) > sd*a.coef {
			a.logger.Info("Found anomaly", slog.Float64("value", freq))
			err = pgdb.AddEntry(response.SessionId, freq, a.logger, response.GetTime().AsTime())
			if err != nil {
				a.logger.Error("Error during insert:", err.Error(), "")
			}
		}
	}
}
