package analyse

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

const (
	minMean = -10
	maxMean = 10
	stdLow  = 0.3
	stdHigh = 1.5
	maxVals = 5000
	tol     = 0.05
)

func areEqual(a, b float64) bool {
	return math.Abs(a-b) < tol
}

func TestMakePrediction(t *testing.T) {
	for j := 0; j < 10; j++ {

		expectedMean := rand.Float64()*(maxMean-minMean) + minMean
		expectedSD := rand.Float64()*(stdHigh-stdLow) + stdLow
		sum := 0.0

		data := make([]float64, maxVals)
		for i := 0; i < maxVals; i++ {
			freq := rand.NormFloat64()*expectedSD + expectedMean
			data[i] = freq
			sum += freq
		}

		a := Analyser{
			maxMessages: maxVals,
			logger:      nil,
			coef:        2.0,
		}

		gotMean, gotSD := a.MakePrediction(data, sum)

		fmt.Printf("test #%d\n", j+1)
		// Verify the results
		if !areEqual(gotMean, expectedMean) {
			t.Errorf("expected mean %f, got %f", expectedMean, gotMean)
		}

		if !areEqual(gotSD, expectedSD) {
			t.Errorf("expected standard deviation %f, got %f", expectedSD, gotSD)
		}
	}
}
