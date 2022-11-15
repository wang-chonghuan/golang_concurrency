package main

import (
	"fmt"
	"runtime"
	"time"
)

/*
*

	type slice struct {
	    array unsafe.Pointer
	    len   int
	    cap   int
	}

type Pointer *ArbitraryType
*/

func callAlgoMatMultiply(mA M, mB M, algo func(M, M) M, title string) M {
	start1 := time.Now()
	retM := algo(mA, mB)
	elapsed := time.Since(start1)
	fmt.Println(title, " elapsed: ", elapsed.Milliseconds(), "ms")
	return retM
}

func test_cannonZero() {
	runtime.GOMAXPROCS(16)
	mA := CreateRandomMatrix(10, 7, 10, 1)
	mB := CreateRandomMatrix(7, 9, 10, 1)
	callAlgoMatMultiply(mA, mB, MultiplyStandard, "MultiplyStandard").printMat("MultiplyStandard")
	callAlgoMatMultiply(mA, mB, MultiplyCannonParallelAnyDefault, "MultiplyCannonParallelAnyDefault").printMat("MultiplyCannonParallelAnyDefault")
}

func test_All() {
	//runtime.GOMAXPROCS(GlobalSqrRootOfNumCores * GlobalSqrRootOfNumCores)
	mA := CreateRandomMatrix(4096, 4096, 1000, 1)
	//mA.printMat("mA")
	mB := CreateRandomMatrix(4096, 4096, 1000, 2)
	//mB.printMat("mB")
	r1 := callAlgoMatMultiply(mA, mB, MultiplySUMMA, "MultiplySUMMA")
	r2 := callAlgoMatMultiply(mA, mB, MultiplyCannonParallelAnyDefault, "MultiplyCannonParallelAnyDefault")
	r3 := callAlgoMatMultiply(mA, mB, MultiplyStandardParallel, "MultiplyStandardParallel")
	if r1.isEqual(&r2) && r2.isEqual(&r3) {
		fmt.Println("all correct")
	}
}

func test_strassen() {
	mA := CreateRandomMatrix(6, 6, 10, 1)
	mA.printMat("mA")
	mB := CreateRandomMatrix(6, 6, 10, 2)
	mB.printMat("mB")
	mStrassen := MultiplyStrassen(mA, mB)
	mStrassen.printMat("mStrassen")
}
