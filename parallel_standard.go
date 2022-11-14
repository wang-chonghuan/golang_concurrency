package main

import "sync"

func sumRowCol(m1 M, m2 M, irm1 int, icm2 int, n int) int {
	var ret int
	for i := 0; i < n; i++ {
		ret += m1[irm1][i] * m2[i][icm2]
	}
	return ret
}

func calcMultiResultRow(m1 M, m2 M, irm1 int, resultRow *([]int), wg1 *sync.WaitGroup) {
	defer wg1.Done()
	_, ncm1 := m1.lens()
	_, ncm2 := m2.lens()
	ret := make([]int, ncm2)
	for icm2 := 0; icm2 < ncm2; icm2++ {
		ret[icm2] = sumRowCol(m1, m2, irm1, icm2, ncm1)
	}
	*resultRow = ret
}

func MultiplyStandard(m1 M, m2 M) (m3 M) {
	if !m1.isMultiplyCompatible(m2) {
		m1.printMat("m1")
		m2.printMat("m2")
		panic("MultiplyStandard NOT isMultiplyCompatible")
	}
	nrm1, ncm1 := m1.lens()
	_, ncm2 := m2.lens()
	m3 = CreateZeroMatrix(nrm1, ncm2)
	for irm1 := 0; irm1 < nrm1; irm1++ {
		for icm2 := 0; icm2 < ncm2; icm2++ {
			m3[irm1][icm2] = sumRowCol(m1, m2, irm1, icm2, ncm1)
		}
	}
	return
}

func MultiplyStandardParallel(m1 M, m2 M) (result M) {
	if !m1.isMultiplyCompatible(m2) {
		return
	}
	nrm1, _ := m1.lens()
	_, ncm2 := m2.lens()
	result = CreateZeroMatrix(nrm1, ncm2)
	var wg1 sync.WaitGroup
	for irm1 := 0; irm1 < nrm1; irm1++ {
		wg1.Add(1)
		go calcMultiResultRow(m1, m2, irm1, &(result[irm1]), &wg1)
	}
	wg1.Wait()
	return
}
