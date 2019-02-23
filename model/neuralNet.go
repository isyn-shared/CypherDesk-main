package model

import (
	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

type ClassificationNetConfiguration struct {
	InputNeurons   int
	OutputNeurons  int
	HiddenNeurons  int
	NumEpochs      int
	LearningRate   float64
	ActivationFunc string
}

type ClassificationNeuralNet struct {
	config                 *ClassificationNetConfiguration
	wHidden                *mat.Dense
	bHidden                *mat.Dense
	wOut                   *mat.Dense
	bOut                   *mat.Dense
	hiddenLayerActivations *mat.Dense
	outputLayerActivations *mat.Dense
	hiddenLayerError       *mat.Dense
	outputLayerError       *mat.Dense
	// networkError *mat.Dense
}

type classificationNetApplies map[string]func(_, _ int, val float64) float64

var (
	cna = make(classificationNetApplies)
)

func PrepareClassConfig() {
	cna["sigmoid"] = func(_, _ int, val float64) float64 { return sigmoid(val) }
	cna["sigmoidDer"] = func(_, _ int, val float64) float64 { return derivativeSigmoid(val) }
}

func newClassificationNeuralNet(nc *ClassificationNetConfiguration) *ClassificationNeuralNet {
	return &ClassificationNeuralNet{config: nc}
}

func (nn *ClassificationNeuralNet) init() {
	nn.hiddenLayerActivations = new(mat.Dense)
	nn.outputLayerActivations = new(mat.Dense)
	nn.hiddenLayerError = new(mat.Dense)
	nn.outputLayerError = new(mat.Dense)
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

func (nn *ClassificationNeuralNet) feedForward(x *mat.Dense) *mat.Dense {
	xRows, xCols := x.Dims()
	if xRows != 1 || xCols != nn.config.InputNeurons {
		panic(fmt.Sprintf("Feedforward error: invalid input data! Input vector has dims [%d, %d], but [1, %d] required",
			xRows, xCols, nn.config.InputNeurons))
	}

	nn.hiddenLayerActivations.Mul(x, nn.wHidden)
	nn.hiddenLayerActivations.Apply(func(_, col int, val float64) float64 {
		return val + nn.bHidden.At(0, col)
	}, nn.hiddenLayerActivations)
	nn.hiddenLayerActivations.Apply(cna[nn.config.ActivationFunc], nn.hiddenLayerActivations)

	nn.outputLayerActivations.Mul(nn.hiddenLayerActivations, nn.wOut)
	nn.outputLayerActivations.Apply(func(_, col int, val float64) float64 {
		return val + nn.bOut.At(0, col)
	}, nn.outputLayerActivations)
	nn.outputLayerActivations.Apply(cna[nn.config.ActivationFunc], nn.outputLayerActivations)

	return nn.outputLayerActivations
}

func (nn *ClassificationNeuralNet) calcErrors(y *mat.Dense) {
	yRows, yCols := y.Dims()
	if yRows != 1 || yCols != nn.config.OutputNeurons {
		panic(fmt.Sprintf("Error when calculating net error: invalid Y! Vector has dims [%d %d], but [1, %d] required",
			yRows, yCols, nn.config.OutputNeurons))
	}
	derApply := cna[nn.config.ActivationFunc+"Der"]
	netError := new(mat.Dense)
	netError.Sub(y, nn.outputLayerActivations)

	slopeOutputLayer := new(mat.Dense)
	slopeHiddenLayer := new(mat.Dense)
	slopeOutputLayer.Apply(derApply, nn.outputLayerActivations)
	slopeHiddenLayer.Apply(derApply, nn.hiddenLayerActivations)

	nn.outputLayerError.MulElem(netError, slopeOutputLayer)
	nn.hiddenLayerError.Mul(nn.outputLayerError, nn.wOut.T())
	nn.hiddenLayerError.MulElem(nn.hiddenLayerError, slopeHiddenLayer)
}

// TODO: Need to understand how this function works
// For scale: https://godoc.org/gonum.org/v1/gonum/mat#Cholesky.Scale
func (nn *ClassificationNeuralNet) backPropagation(x *mat.Dense) {
	wOutAdj := new(mat.Dense)
	wOutAdj.Mul(nn.hiddenLayerActivations.T(), nn.outputLayerError)
	wOutAdj.Scale(nn.config.LearningRate, wOutAdj)
	nn.wOut.Add(nn.wOut, wOutAdj)

	bOutAdj := sumAlongAxis(0, nn.outputLayerError)
	bOutAdj.Scale(nn.config.LearningRate, bOutAdj)
	nn.bOut.Add(nn.bOut, bOutAdj)

	wHiddenAdj := new(mat.Dense)
	wHiddenAdj.Mul(x.T(), nn.hiddenLayerError)
	wHiddenAdj.Scale(nn.config.LearningRate, wHiddenAdj)
	nn.wHidden.Add(nn.wHidden, wHiddenAdj)

	bHiddenAdj := sumAlongAxis(1, nn.hiddenLayerError)
	bHiddenAdj.Scale(nn.config.LearningRate, bHiddenAdj)
	nn.bHidden.Add(nn.bHidden, bHiddenAdj)
}

func (nn *ClassificationNeuralNet) train(x, y *mat.Dense) {
	for i := 0; i < nn.config.NumEpochs; i++ {
		fmt.Printf("Epoch number %d was started", i + 1)
		nn.feedForward(x)
		nn.calcErrors(y)
		fmt.Println("Sigma output layer error:")
		printDense(nn.outputLayerError)
		nn.backPropagation(x)
	}
}

func Debug() {
	PrepareClassConfig()
	nnConf := &ClassificationNetConfiguration{
		InputNeurons:   10,
		HiddenNeurons:  7,
		OutputNeurons:  3,
		ActivationFunc: "sigmoid",
	}
	nn := newClassificationNeuralNet(nnConf)
	nn.init()

	x := make([]float64, nn.config.InputNeurons)
	y := make([]float64, nn.config.OutputNeurons)
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range x {
		x[i] = randGen.Float64()
	}
	for i := range y {
		y[i] = randGen.Float64()
	}

	xConv := ConvertInputValues(x, 1, nn.config.InputNeurons)
	yConv := ConvertInputValues(y, 1, nn.config.OutputNeurons)
	out := nn.feedForward(xConv)

	fmt.Println("OUTPUT: ")
	printDense(out)

	nn.calcErrors(yConv)

	fmt.Println("ERRORS")
	fmt.Println("outputLayerError: ")
	printDense(nn.outputLayerError)
	fmt.Println("hiddenLayerError: ")
	printDense(nn.hiddenLayerError)
}
