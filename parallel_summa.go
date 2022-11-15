package main

import (
	"fmt"
	"sync"
)

type StripDim struct {
	addR1, addC2 int
	wStrip       int
}

func ForCompetition22117989(m1 [][]int, m2 [][]int) [][]int {
	return [][]int(MultiplySUMMA(M(m1), M(m2)))
}

func MultiplySUMMA(mA M, mB M) M {
	sd := StripDim{}
	optimizeMat(&mA, &mB, &sd)
	smA := CreateStripMat(mA, true, sd.wStrip, mA.nc(), mA.nr()/sd.wStrip)
	smB := CreateStripMat(mB, false, sd.wStrip, mB.nr(), mB.nc()/sd.wStrip)
	result := smA.multiply(&smB)
	recoverMat(&mA, &mB, &result, &sd)
	return result
}

func recoverMat(m1 *M, m2 *M, mr *M, sd *StripDim) {
	if sd.addR1 == 1 {
		m1.removeLastRow()
		mr.removeLastRow()
	}
	if sd.addC2 == 1 {
		m2.removeLastCol()
		mr.removeLastCol()
	}
}

func optimizeMat(m1 *M, m2 *M, sd *StripDim) {
	if m1.nr()*m1.nc()*m2.nr()*m2.nc() == 0 || m1.nc() != m2.nr() {
		panic("optimizeMat: dims not compatible")
	}
	if m1.nr()%2 != 0 {
		sd.addR1 += 1
		m1.addZeroRow()
	}
	if m2.nc()%2 != 0 {
		sd.addC2 += 1
		m2.addZeroCol()
	}
	sd.wStrip = 1
	tempWidth := 1
	var nStrip1, nStrip2 int
	ratio := float64(float64(max(m1.nr()*m2.nr(), 4096*4096)) / float64(4096*4096))
	maxG := int(MaxGoroutines * ratio)
	for {
		if tempWidth > min(m1.nr(), m2.nc()) {
			break
		}
		if m1.nr()%tempWidth == 0 && m2.nc()%tempWidth == 0 {
			sd.wStrip = tempWidth
			nStrip1 = m1.nr() / tempWidth
			nStrip2 = m2.nc() / tempWidth
			nThread := nStrip1 * nStrip2 // tempWidth 增大， nThread减少
			if nThread > maxG {
				tempWidth++
				continue
			} else {
				break
			}
		} else {
			tempWidth++
			continue
		}
	}
	fmt.Println(fmt.Sprintf("%+v nStrip1 %v nStrip2 %v r1 %v c1 %v r2 %v c2 %v",
		sd, nStrip1, nStrip2, m1.nr(), m1.nc(), m2.nr(), m2.nc()))
	return
}

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
	//isRow      bool    // are these strips listed as rows; if false, its cols
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
