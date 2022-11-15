package main

// Chonghuan Wang 22117989
func MultiplyCannonParallelSquare(s1 M, s2 M, dimByBlock int, dimOfBlock int) M {
	bmA := CreateBlockSqrMat(s1, dimOfBlock, dimByBlock)
	bmB := CreateBlockSqrMat(s2, dimOfBlock, dimByBlock)
	bmC := CreateEmptyBlockMat(dimByBlock, dimOfBlock)
	channles := CreateChannels(dimByBlock)
	cluster := CreatCluster(dimByBlock, dimOfBlock, bmA, bmB, channles)
	cluster.goProcs()
	for i := 0; i < dimByBlock; i++ {
		cluster.broadcastNextBlocksToProcs(bmA, bmB, i)
	}
	cluster.closeProcs()
	cluster.joinBlocks(dimOfBlock, &bmC.m)
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
