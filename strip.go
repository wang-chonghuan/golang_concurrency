package main

import "sync"

// strip.go

type Strip struct {
	wStrip        int // number of elements in a substrip
	lStrip        int // number of substrips in a strip
	substripArray M   // substrip listed in a strip
}

func (o *Strip) multipyByOuterProduct(right *Strip) M {
	outerProduct := CreateZeroMatrix(o.wStrip, right.wStrip)
	for iSubstrip := 0; iSubstrip < o.lStrip; iSubstrip++ {
		for iLeft := 0; iLeft < o.wStrip; iLeft++ {
			for iRight := 0; iRight < right.wStrip; iRight++ {
				outerProduct[iLeft][iRight] += o.substripArray[iSubstrip][iLeft] * right.substripArray[iSubstrip][iRight]
			}
		}
	}
	return outerProduct
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
	nStrip     int     // how many strips in the matrix
	stripArray []Strip // strip listed in a matrix
}

func (o *StripMat) multiply(right *StripMat) M {
	nrRes := o.wStrip * o.nStrip
	ncRes := right.wStrip * right.nStrip
	result := CreateZeroMatrix(nrRes, ncRes)
	var wg sync.WaitGroup
	for iLeft := 0; iLeft < o.nStrip; iLeft++ {
		for iRight := 0; iRight < right.nStrip; iRight++ {
			wg.Add(1)
			go stripMultiply(&(o.stripArray[iLeft]), &(right.stripArray[iRight]),
				&result, iLeft, iRight, &wg)
		}
	}
	wg.Wait()
	return result
}

func stripMultiply(leftStrip *Strip, rightStrip *Strip, result *M, iLeft int, iRight int, wg *sync.WaitGroup) {
	defer wg.Done()
	blockM := leftStrip.multipyByOuterProduct(rightStrip)
	nrBm, _ := blockM.lens()
	joinMat(*result, blockM, iRight*nrBm, iRight*nrBm+nrBm, iLeft*nrBm, iLeft*nrBm+nrBm)
}

func CreateStripMat(srcMat M, isRow bool, wStrip int, lStrip int, nStrip int) StripMat {
	nr, nc := srcMat.lens()
	sm := StripMat{wStrip: wStrip, nStrip: nStrip, stripArray: make([]Strip, nStrip)}
	if !((isRow && wStrip*nStrip == nr && lStrip == nc) || (!isRow && wStrip*nStrip == nc && lStrip == nr)) {
		panic("CreateStripMat size not compatible")
	}
	for iStrip := 0; iStrip < nStrip; iStrip++ {
		sm.stripArray[iStrip] = CreateEmptyStrip(wStrip, lStrip)
		sm.stripArray[iStrip].copyMatToStrip(&srcMat, isRow, wStrip*iStrip, wStrip*iStrip+wStrip)
	}
	return sm
}
