package main

import "fmt"

type Strip struct {
	wStrip        int  // number of elements in a substrip
	lStrip        int  // number of substrips in a strip
	substripArray M    // substrip listed in a strip
	isRow         bool // !!!if the strip is listed as rows(when isRow is true), its substrips are listed as cols!!!
}

func (o *Strip) multipyByOuterProduct(right *Strip) M {
	result := CreateZeroMatrix(o.wStrip, right.wStrip)
	outerProduct := CreateZeroMatrix(o.wStrip, right.wStrip)
	for iSubstrip := 0; iSubstrip < o.lStrip; iSubstrip++ {
		for iLeft := 0; iLeft < o.wStrip; iLeft++ {
			for iRight := 0; iRight < right.wStrip; iRight++ {
				outerProduct[iLeft][iRight] = o.substripArray[iSubstrip][iLeft] * right.substripArray[iSubstrip][iRight]
			}
		}
		outerProduct.printMat(fmt.Sprintf("outerProduct %v", iSubstrip))
		result = addMat(result, outerProduct)
	}
	return result
}

func CreateEmptyStrip(wStrip int, lStrip int) Strip {
	s := Strip{}
	s.wStrip = wStrip
	s.lStrip = lStrip
	s.substripArray = make([][]int, lStrip)
	for i := 0; i < lStrip; i++ {
		s.substripArray[i] = make([]int, wStrip)
	}
	return s
}

// 此处可化简
func (o *Strip) copyMatToStrip(srcMat *M, isRow bool, iBegin int, iEnd int) {
	if isRow {
		for ir := iBegin; ir < iEnd; ir++ {
			for ic := 0; ic < o.lStrip; ic++ {
				o.substripArray[ic][ir-iBegin] = (*srcMat)[ir][ic]
			}
		}
	} else {
		for ic := iBegin; ic < iEnd; ic++ {
			for ir := 0; ir < o.lStrip; ir++ {
				o.substripArray[ir][ic-iBegin] = (*srcMat)[ir][ic]
			}
		}
	}
}

type StripMat struct {
	wStrip     int     // number of elements in a substrip
	lStrip     int     // number of substrips in a strip
	nStrip     int     // how many strips in the matrix
	stripArray []Strip // strip listed in a matrix
	isRow      bool    // are these strips listed as rows; if false, its cols
}

func (o *StripMat) multiply(right *StripMat) M {
	nrRes := o.wStrip * o.nStrip
	ncRes := right.wStrip * right.nStrip
	result := CreateZeroMatrix(nrRes, ncRes)
	for iLeft := 0; iLeft < o.nStrip; iLeft++ {
		for iRight := 0; iRight < right.nStrip; iRight++ {
			blockM := o.stripArray[iLeft].multipyByOuterProduct(&(right.stripArray[iRight]))
			nrBm, ncBm := blockM.lens()
			if nrBm != ncBm {
				panic("(o *StripMat) multiply nrBm != ncBm")
			}
			joinMat(result, blockM, iRight*nrBm, iRight*nrBm+nrBm, iLeft*nrBm, iLeft*nrBm+nrBm)
		}
	}
	return result
}

func CreateStripMat(srcMat M, isRow bool, wStrip int, lStrip int, nStrip int) StripMat {
	nr, nc := srcMat.lens()
	sm := StripMat{isRow: isRow, wStrip: wStrip, lStrip: lStrip, nStrip: nStrip, stripArray: make([]Strip, nStrip)}
	if !((isRow && wStrip*nStrip == nr && lStrip == nc) || (!isRow && wStrip*nStrip == nc && lStrip == nr)) {
		panic("CreateStripMat size not compatible")
	}
	for iStrip := 0; iStrip < nStrip; iStrip++ {
		sm.stripArray[iStrip] = CreateEmptyStrip(wStrip, lStrip)
		sm.stripArray[iStrip].copyMatToStrip(&srcMat, isRow, wStrip*iStrip, wStrip*iStrip+wStrip)
	}
	return sm
}
