package model

import (
	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

type neuralNetConf struct {
	InputNeurons  int
	OutputNeurons int
	HiddenNeurons int
	NumEpochs     int
	LearningRate  float64
}

type neuralNet struct {
	config  *neuralNetConf
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOut    *mat.Dense
	bOut    *mat.Dense
}

func newNeuralNet(nc *neuralNetConf) *neuralNet {
	return &neuralNet{config: nc}
}

func (nn *neuralNet) init() {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	nn.wHidden = mat.NewDense(nn.config.InputNeurons, nn.config.HiddenNeurons, nil)
	nn.bHidden = mat.NewDense(1, nn.config.HiddenNeurons, nil)
	nn.wOut = mat.NewDense(nn.config.HiddenNeurons, nn.config.OutputNeurons, nil)
	nn.bOut = mat.NewDense(1, nn.config.OutputNeurons, nil)

	wHiddenRaw := nn.wHidden.RawMatrix().Data
	bHiddenRaw := nn.bHidden.RawMatrix().Data
	wOutRaw := nn.wOut.RawMatrix().Data
	bOutRaw := nn.bOut.RawMatrix().Data

	for _, param := range [][]float64{
		wHiddenRaw,
		bHiddenRaw,
		wOutRaw,
		bOutRaw,
	} {
		for i := range param {
			param[i] = randGen.Float64()
		}
	}
}

func (nn *neuralNet) feedForward(x *mat.Dense) *mat.Dense {
	xRows, xCols := x.Dims()
	if xRows != 1 || xCols != nn.config.InputNeurons {
		panic(fmt.Sprintf("Feedforward error: invalid input data! Input vector has dims [%d, %d], but [1, %d] required",
			xRows, xCols, nn.config.InputNeurons))
	}
	sigApply := func(_, _ int, val float64) float64 { return sigmoid(val) }
	// derSigAderSigApply	pply := func(_, _ int, val float64) float64 { return derivativeSigmoid(val) }

	hiddenVals := new(mat.Dense)
	hiddenVals.Mul(x, nn.wHidden)
	hiddenVals.Apply(func(_, col int, val float64) float64 {
		return val + nn.bHidden.At(0, col)
	}, hiddenVals)
	hiddenVals.Apply(sigApply, hiddenVals)

	outputVals := new(mat.Dense)
	outputVals.Mul(hiddenVals, nn.wOut)
	outputVals.Apply(func(_, col int, val float64) float64 {
		return val + nn.bOut.At(0, col)
	}, outputVals)
	outputVals.Apply(sigApply, outputVals)

	return outputVals
}

func (nn *neuralNet) calcNetError() {

}

func (nn *neuralNet) backPropogation() {

}

func Debug() {
	nnConf := &neuralNetConf{
		InputNeurons:  10,
		HiddenNeurons: 7,
		OutputNeurons: 3,
	}
	nn := newNeuralNet(nnConf)
	nn.init()

	x := make([]float64, nn.config.InputNeurons)
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range x {
		x[i] = randGen.Float64()
	}

	xConv := ConvertInputValues(x, 1, nn.config.InputNeurons)
	y := nn.feedForward(xConv)

	printDense(y)
}
