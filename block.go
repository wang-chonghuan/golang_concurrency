package main

import "fmt"

// block matrix, which is square matrix
type BlockMat struct {
	bm         [][]M // block matrix
	m          M     // flat matrix
	dim        int   // dim by blocks
	dimOfBlock int   // dim of each block
}

func CreateEmptyBlockMat(dimByBlock int, dimOfBlock int) BlockMat {
	bmC := BlockMat{}
	bmC.init(dimByBlock, dimOfBlock)
	return bmC
}

func CreateBlockSqrMat(srcM M, dimOfBlock int, dimByBlock int) BlockMat {
	//nr, nc := srcM.lens()
	//fmt.Println("CreateBlockSqrMat: dimOfBlock%v dimByBlock%v", dimOfBlock, dimByBlock, nr, nc)
	var bsm BlockMat
	bsm.init(dimByBlock, dimOfBlock)
	for ir := 0; ir < bsm.dim; ir++ {
		for ic := 0; ic < bsm.dim; ic++ {
			bsm.bm[ir][ic] = SplitMat(srcM, dimOfBlock, dimOfBlock,
				ic*dimOfBlock, ic*dimOfBlock+dimOfBlock, ir*dimOfBlock, ir*dimOfBlock+dimOfBlock)
		}
	}
	return bsm
}

func (o *BlockMat) init(dimByBlock int, dimOfBlock int) {
	o.dim = dimByBlock
	o.dimOfBlock = dimOfBlock
	o.m = CreateZeroMatrix(dimOfBlock*dimByBlock, dimOfBlock*dimByBlock)
	o.bm = make([][]M, dimByBlock)
	for iRow := 0; iRow < dimByBlock; iRow++ {
		o.bm[iRow] = make([]M, dimByBlock)
		for iBlock := 0; iBlock < dimByBlock; iBlock++ {
			o.bm[iRow][iBlock] = CreateZeroMatrix(dimOfBlock, dimOfBlock)
		}
	}
}

func (o *BlockMat) getBlock(ir int, ic int) M {
	return o.bm[ir][ic]
}

func (o *BlockMat) setOneBlock(ir int, ic int, m M) {
	m.panicIfWrongDim(o.dimOfBlock, o.dimOfBlock)
	o.bm[ir][ic] = m
}

func (o *BlockMat) setBlockRow(ir int, blockRow []M) {
	for ic := 0; ic < o.dim; ic++ {
		o.setOneBlock(ir, ic, blockRow[ic])
	}
}

func (o *BlockMat) setBlockCol(ic int, blockCol []M) {
	for ir := 0; ir < o.dim; ir++ {
		o.setOneBlock(ir, ic, blockCol[ir])
	}
}

func (o *BlockMat) rotateRowsUpOneStep() {
	if o.dim <= 1 {
		return
	} else {
		overflowRow := make([]M, o.dim)
		copy(overflowRow, o.bm[0])
		for ir := 1; ir <= o.dim-1; ir++ {
			o.setBlockRow(ir-1, o.bm[ir])
		}
		o.setBlockRow(o.dim-1, overflowRow)
	}
}

func (o *BlockMat) rotateColsLeftOneStep() {
	if o.dim <= 1 {
		return
	} else {
		for ir := 0; ir < o.dim; ir++ {
			rotateLeftOneStep(&(o.bm[ir]))
		}
	}
}

func rotateLeftOneStep(row *[]M) {
	n := len(*row)
	if n <= 1 {
		return
	} else {
		overflowM := (*row)[0]
		for i := 1; i <= n-1; i++ {
			(*row)[i-1] = (*row)[i]
		}
		(*row)[n-1] = overflowM
	}
	return
}

func (o *BlockMat) printMat(title string) {
	fmt.Println(title + " begin")
	for iR := 0; iR < o.dim; iR++ {
		for iC := 0; iC < o.dim; iC++ {
			o.bm[iR][iC].printMat(fmt.Sprintf("block %v %v", iR, iC))
		}
	}
	fmt.Println(title + " end")
}
