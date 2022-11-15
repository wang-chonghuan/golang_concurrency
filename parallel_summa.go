package main

import "fmt"

type StripDim struct {
	addR1, addC2 int
	wStrip       int
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
	maxG := MaxGoroutines * m1.nr() * m2.nr() / (4096 * 4096)
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
