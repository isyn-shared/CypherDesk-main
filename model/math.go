package model

import (
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func derivativeSigmoid(x float64) float64 {
	return x * (1.0 - x)
}

func sumAlongAxis(axis int, m *mat.Dense) *mat.Dense {
	rows, cols := m.Dims()
	var res *mat.Dense

	switch axis {
	case 0:
		data := make([]float64, cols)
		for i := 0; i < cols; i++ {
			col := mat.Col(nil, i, m)
			data[i] = floats.Sum(col)
		}
		res = mat.NewDense(1, cols, data)
	case 1:
		data := make([]float64, rows)
		for i := 0; i < rows; i++ {
			row := mat.Row(nil, i, m)
			data[i] = floats.Sum(row)
		}
		res = mat.NewDense(rows, 1, data)
	default:
		panic("invalid axis, must be 0 or 1")
	}

	return res
}
