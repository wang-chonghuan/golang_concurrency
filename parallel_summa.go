package main

import "fmt"

func test_summa() {
	mA := CreateRandomMatrix(8, 6, 10, 1)
	mB := CreateRandomMatrix(6, 12, 10, 2)
	mA.printMat("mA")
	mB.printMat("mB")

	smA := CreateStripMat(mA, true, 4, 6, 2)
	for i := 0; i < smA.nStrip; i++ {
		smA.stripArray[i].substripArray.printMat(fmt.Sprintf("SMA substripArray %v", i))
	}
	smB := CreateStripMat(mB, false, 4, 6, 3)
	for i := 0; i < smB.nStrip; i++ {
		smB.stripArray[i].substripArray.printMat(fmt.Sprintf("SMB substripArray %v", i))
	}

}
