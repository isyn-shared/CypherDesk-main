package model

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func printDense(d *mat.Dense) {
	rows, cols := d.Dims()
	fmt.Printf("Dimensions: %d, %d\n", rows, cols)
	data := d.RawMatrix().Data
	for i, val := range data {
		if (i+1)%cols == 0 && i > 0 {
			fmt.Printf("%f\n", val)
		} else {
			fmt.Printf("%f ", val)
		}
	}
}

func (nn *neuralNet) Print() {
	fmt.Println("******** HIDDEN LAYOUT *********")
	printDense(nn.wHidden)
	fmt.Println("BIOS")
	printDense(nn.bHidden)

	fmt.Println("********* OUTPUT LAYOUT ********")
	printDense(nn.wOut)
	fmt.Println("BIOS")
	printDense(nn.bOut)
}

func ConvertInputValues(vals []float64, n, m int) *mat.Dense {
	return mat.NewDense(n, m, vals)
}
