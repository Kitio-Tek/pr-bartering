package failuresimulation

import (
	configextractor "bartering/config-extractor"
	"bartering/utils"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

func ExtractFailureModel(config configextractor.Config) (func(float64, float64) float64, error) {

	/*
		Given config, extract the node's failure model
		(failure model is the probability law for session length)
	*/

	switch config.FailureModel {
	case "weibull":
		return DrawNumberWeibull, nil
	case "lognormal":
		return DrawNumberLognormal, nil
	default:
		return func(float64, float64) float64 { return 0.0 }, errors.New("failure model not recognized")
	}
}

func ExtractConnectivityFactor(config configextractor.Config) (float64, error) {

	/*
		Given config, extract node profile
		(node profile defines proportion of time where node is up or down)
	*/

	switch config.NodeProfile {
	case "peer":
		return 0.5, nil
	case "benefactor":
		return 0.7, nil
	case "peeper":
		return 0.3, nil
	default:
		return 0.0, errors.New("node profile not recognized ; should be benefactor, peer or peeper")
	}
}

func DrawNumberWeibull(shape float64, scale float64) float64 {

	/*
		Draw session length according to weibull law
	*/

	weibullDist := distuv.Weibull{
		K:      shape,
		Lambda: scale,
	}
	return weibullDist.Rand()
}

func DrawNumberLognormal(Mu float64, Sigma float64) float64 {

	/*
		Draw session length according to lognormal law
	*/

	logNormalDist := distuv.LogNormal{
		Mu:    Mu,
		Sigma: Sigma,
	}
	return logNormalDist.Rand()
}

func computeDowntimeFromSessionLength(connectivityFactor float64, sessionLength float64) float64 {

	/*
		Given connectivity factor, we compute downtime to ensure over one cycle of up-down, connectivity factor is respected
	*/

	return ((1 - connectivityFactor) / connectivityFactor) * sessionLength

}

func stopNode(mutex *sync.Mutex, downTime float64) {

	/*
		Acquires a mutex that blocks peers listener and simulates downtime for a node
	*/

	fmt.Println("Stopping node for ", downTime)
	mutex.Lock()
	fmt.Println("Mutex locked, no communication now")
	time.Sleep(time.Duration(downTime) * time.Second)
	mutex.Unlock()
	fmt.Println("Mutex unlocked")
}

func Failure(config configextractor.Config, shape float64, scale float64, mutex *sync.Mutex) {

	/*
		Failure func, given config, probability law parameters and a mutex, simulates failure
	*/

	// sessionLengthDraw, err := ExtractFailureModel(config)

	// utils.ErrorHandler(err)

	connectivityFactor, err := ExtractConnectivityFactor(config)

	utils.ErrorHandler(err)

	for {

		// sessionLength := sessionLengthDraw(shape, scale)
		sessionLength := 10.0

		downTime := computeDowntimeFromSessionLength(connectivityFactor, sessionLength)

		sessionLengthStr := strconv.FormatFloat(sessionLength, 'f', -1, 64)
		// downTimeStr := strconv.FormatFloat(downTime, 'f', -1, 64)
		fmt.Println("Staying up for ", sessionLengthStr)
		time.Sleep(time.Duration(sessionLength) * time.Second)
		stopNode(mutex, downTime)
	}

}
