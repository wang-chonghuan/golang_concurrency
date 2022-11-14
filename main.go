package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	// we are seeding the rand variable with present time
	// so that we would get different output each time
}

func main() {
	//test_All()
	test_summa()
}

func main1() {
	var col = 4096
	var row = 4096
	var randMatrixA [][]int
	randMatrixA = genMat(row, col)
	fmt.Println(randMatrixA)
	var randMatrixB [][]int
	randMatrixB = genMat(col, row)
	fmt.Println(randMatrixB)
	fmt.Println("The Go Result of Matrix Multiplication = ")
	start := time.Now()
	//c := doCalc(randMatrixA, randMatrixB)
	//c := MultiplyStandardParallel(randMatrixA, randMatrixB)
	c := MultiplyCannonParallelAnyDefault(randMatrixA, randMatrixB)
	//c := MultiplyStrassen(randMatrixA, randMatrixB)
	elapsed := time.Since(start)

	fmt.Println(c)

	fmt.Printf("Time taken to calculate %s \n", elapsed)
}
func genMat(row int, col int) [][]int {
	nM := make([][]int, col)
	for i := 0; i < col; i++ {
		nM[i] = make([]int, row)
		// we are creating a slice which can hold type int
	}
	generateNums(nM)
	return nM
}
func generateNums(randMatrix [][]int) {
	for i, innerArray := range randMatrix {
		for j := range innerArray {
			randMatrix[i][j] = rand.Intn(100)
			//looping over each element of array and assigning it a random variable
		}
	}
}
func rowCount(inM [][]int) int {
	return (len(inM))
}
func colCount(inM [][]int) int {
	return (len(inM[0]))
}
func doCalc(inA [][]int, inB [][]int) [][]int {
	var i, j int
	m := rowCount(inA) // number of rows the first matrix
	//   n := colCount(inA)     // number of columns the first matrix
	p := rowCount(inB) // number of rows the second matrix
	q := colCount(inB) // number of columns the second matrix
	k := 0
	total := 0
	var nM [][]int
	nM = genMat(m, q)
	for i = 0; i < m; i++ {
		for j = 0; j < q; j++ {
			for k = 0; k < p; k++ {
				total = total + inA[i][k]*inB[k][j]
				//      fmt.Print("(", inA[i][k], " * ", inB[k][j], ") + ")
			}
			//          fmt.Println("giving", total)
			nM[i][j] = total
			total = 0
		}
		fmt.Println()
	}
	return nM
}
