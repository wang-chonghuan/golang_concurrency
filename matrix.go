package main

import (
	"fmt"
	"math/rand"
)

type M [][]int

func CreateZeroMatrix(rows int, cols int) M {
	m := make([][]int, rows)
	for i := 0; i < rows; i++ {
		m[i] = make([]int, cols)
	}
	return m
}

// rand.Seed(time.Now().UnixNano())
func CreateRandomMatrix(rows int, cols int, upBoundry int, seed int) M {
	m := make([][]int, rows)
	rand.Seed(int64(seed))
	for ir := 0; ir < rows; ir++ {
		m[ir] = make([]int, cols)
		for ic := 0; ic < cols; ic++ {
			m[ir][ic] = int(rand.Intn(upBoundry))
		}
	}
	return m
}

func (o *M) init(rows int, cols int) {
	m := make([][]int, cols)
	for i := 0; i < rows; i++ {
		m[i] = make([]int, cols)
	}
	o = (*M)(&m)
}

func (o *M) nr() int {
	return len(*o)
}

func (o *M) nc() int {
	return len((*o)[0])
}

func (o M) lens() (int, int) {
	rows := len(o)
	if rows == 0 {
		return 0, 0
	} else {
		return rows, len(o[0])
	}
}

func (o M) isMultiplyCompatible(m2 M) bool {
	rows1, cols1 := o.lens()
	rows2, cols2 := m2.lens()
	if cols1 == rows2 && cols1 != 0 {
		return true
	} else {
		fmt.Println(fmt.Sprintf("%v %v %v %v", rows1, cols1, rows2, cols2))
		return false
	}
}

func (o M) printMat(title string) {
	fmt.Println(title + " begin")
	rows, cols := o.lens()
	for iR := 0; iR < rows; iR++ {
		for iC := 0; iC < cols; iC++ {
			fmt.Printf("%v ", o[iR][iC])
		}
		fmt.Println()
	}
	fmt.Println(title + " end")
}

func (o M) panicIfWrongDim(nr int, nc int) {
	rows, cols := o.lens()
	if nr != rows || nc != cols {
		panic(fmt.Sprintf("this matrix(%vx%v) should be a %vx%v matrix", rows, cols, nr, nc))
	}
}

func (o *M) isEqual(right *M) bool {
	nr1, nc1 := o.lens()
	nr2, nc2 := right.lens()
	if nr1 != nr2 || nc1 != nc2 {
		return false
	}
	for ir := 0; ir < nr1; ir++ {
		for ic := 0; ic < nc1; ic++ {
			if (*o)[ir][ic] != (*right)[ir][ic] {
				return false
			}
		}
	}
	return true
}

func (m *M) addZeroCol() {
	for ir := 0; ir < m.nr(); ir++ {
		(*m)[ir] = append((*m)[ir], 0)
	}
}

func (m *M) removeLastCol() {
	for ir := 0; ir < m.nr(); ir++ {
		(*m)[ir] = (*m)[ir][:len((*m)[ir])-1]
	}
}

func (m *M) addZeroRow() {
	*m = append(*m, make([]int, m.nc()))
}

func (m *M) removeLastRow() {
	*m = (*m)[:len(*m)-1]
}

func PrintRowsCols(m1 M, m2 M) {
	rows1, cols1 := m1.lens()
	rows2, cols2 := m2.lens()
	fmt.Printf("rows1 %d cols1 %d rows2 %d cols2 %d, multi %d %d\n",
		rows1, cols1, rows2, cols2, rows1, cols2)
}

func test_addAndSubMat() {
	var m1 M = [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	var m2 M = [][]int{
		{3, 2, 1},
		{6, 5, 4},
	}
	m1.printMat("m1")
	m2.printMat("m2")
	addMat(m1, m2).printMat("add")
	subMat(m1, m2).printMat("sub")

}

func subMat(m1 M, m2 M) M {
	return binaryOptMat(m1, m2, func(a int, b int) int {
		return a - b
	})
}

func addMat(m1 M, m2 M) M {
	return binaryOptMat(m1, m2, func(a int, b int) int {
		return a + b
	})
}

func binaryOptMat(m1 M, m2 M, f func(int, int) int) M {
	nrm1, ncm1 := m1.lens()
	nrm2, ncm2 := m2.lens()
	if nrm1 != nrm2 || ncm1 != ncm2 {
		return M{}
	}
	result := CreateZeroMatrix(nrm1, ncm1)
	for irm1 := 0; irm1 < nrm1; irm1++ {
		for icm1 := 0; icm1 < ncm1; icm1++ {
			result[irm1][icm1] = f(m1[irm1][icm1], m2[irm1][icm1])
		}
	}
	return result
}

func joinMat(dstM M, srcM M, left int, right int, up int, down int) {
	iSrcM := 0
	for iDstM := up; iDstM < down; iDstM++ {
		copy(dstM[iDstM][left:right], srcM[iSrcM])
		iSrcM++
	}
}

// section: beginning included, ending NOT include
func SplitMat(srcM M, nrDst int, ncDst int, left int, right int, up int, down int) M {
	nrSrc, ncSrc := srcM.lens()
	if nrSrc == nrDst && ncSrc == ncDst {
		fmt.Println("SplitMat: no need")
		return srcM
	}
	dstM := CreateZeroMatrix(nrDst, ncDst)
	iDstM := 0
	for iSrcM := up; iSrcM < down; iSrcM++ {
		copy(dstM[iDstM], srcM[iSrcM][left:right])
		iDstM++
	}
	return dstM
}

func MatToSqr(m1 M, m2 M, dim int) (sqr1 M, sqr2 M) {
	nrm1, ncm1 := m1.lens()
	nrm2, ncm2 := m2.lens()
	//fmt.Println(fmt.Sprintf("MatToSqr: nr1 %v nc1 %v nr2 %v nc2 %v dimOfSquare %v", nrm1, ncm1, nrm2, ncm2, dim))
	if (nrm1 == ncm1) && (nrm2 == ncm2) && (ncm1 == nrm2) && (nrm1 == dim) {
		fmt.Println("MatToSqr: no need")
		return m1, m2
	}
	sqr1 = CreateZeroMatrix(dim, dim)
	sqr2 = CreateZeroMatrix(dim, dim)
	fillZero(m1, sqr1)
	fillZero(m2, sqr2)
	return sqr1, sqr2
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func pow(base int, exp int) int {
	if exp == 0 {
		return base
	} else {
		return 2 * pow(base, exp-1)
	}
}

func min2n(a int) (dim int) {
	exp := 0
	for {
		cur2n := pow(2, exp)
		if cur2n < a {
			exp++
			continue
		} else {
			return cur2n
		}
	}
}

func calc2nDim(nr1, nc1, nr2, nc2 int) int {
	if nc1 != nr2 {
		panic("calc2nDim nc1 != nrm2")
	}
	maxN := max(max(nr1, nc1), nc2)
	min2expN := min2n(maxN)
	return min2expN
}

func calcByCores(nr1 int, nc1 int, nr2 int, nc2 int, sqrRootCores int) int {
	if nc1 != nr2 {
		panic("calcByCores nc1 != nr2")
	}
	maxLen := max(nr1, max(nc1, nc2))
	i := 1
	for {
		curLen := i * sqrRootCores
		if curLen >= maxLen {
			return curLen
		} else {
			i++
		}
	}
}
