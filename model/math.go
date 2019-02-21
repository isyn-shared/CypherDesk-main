package model

import "math"

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func derivativeSigmoid(x float64) float64 {
	return x * (1.0 - x)
}