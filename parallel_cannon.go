package main

import (
	"fmt"
	"time"
)

// Chonghuan Wang 22117989
func MultiplyCannonParallelSquare(s1 M, s2 M, dimByBlock int, dimOfBlock int) M {
	t1 := time.Now()
	bmA := CreateBlockSqrMat(s1, dimOfBlock, dimByBlock)
	bmB := CreateBlockSqrMat(s2, dimOfBlock, dimByBlock)
	bmC := CreateEmptyBlockMat(dimByBlock, dimOfBlock)
	fmt.Println(time.Since(t1).Milliseconds())
	channles := CreateChannels(dimByBlock)
	cluster := CreatCluster(dimByBlock, dimOfBlock, bmA, bmB, channles)
	cluster.goProcs()
	t3 := time.Now()
	for i := 0; i < dimByBlock; i++ {
		cluster.broadcastNextBlocksToProcs(bmA, bmB, i)
	}
	fmt.Println(time.Since(t3).Milliseconds())
	t2 := time.Now()
	cluster.closeProcs()
	cluster.joinBlocks(dimOfBlock, &bmC.m)
	fmt.Println(time.Since(t2).Milliseconds())
	return bmC.m
}

func MultiplyCannonParallelAnyDefault(m1 M, m2 M) M {
	return MultiplyCannonParallelAny(m1, m2, GlobalSqrRootOfNumCores)
}

func MultiplyCannonParallelAny(m1 M, m2 M, sqrRootOfNumCores int) M {
	nr1, nc1 := m1.lens()
	nr2, nc2 := m2.lens()
	dimOfZeroMat := calcByCores(nr1, nc1, nr2, nc2, sqrRootOfNumCores)
	s1, s2 := MatToSqr(m1, m2, dimOfZeroMat)
	dimOfBlock := dimOfZeroMat / sqrRootOfNumCores
	dimByBlock := sqrRootOfNumCores
	//fmt.Println(fmt.Sprintf("MultiplyCannonParallelAny: sqrRootOfNumCores %v, dimOfZeroMat %v, dimOfBlock %v, dimByBlock %v", sqrRootOfNumCores, dimOfZeroMat, dimOfBlock, dimByBlock))
	retS := MultiplyCannonParallelSquare(s1, s2, dimByBlock, dimOfBlock)
	retM := SplitMat(retS, nr1, nc2, 0, nc2, 0, nr1)
	return retM
}
