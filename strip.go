package main

type Strip struct {
	wStrip        int // number of elements in a substrip
	lStrip        int // number of substrips in a strip
	substripArray M   // substrip listed in a strip
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
