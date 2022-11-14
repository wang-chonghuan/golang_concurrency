package main

func fillZero(m M, sqr M) {
	for k, v := range m {
		copy(sqr[k], v)
	}
}

// input: 2 squares with same dim
func doMultiplyStrassen(ret *M, s1 M, s2 M) {
	n, _ := s1.lens()
	*ret = CreateZeroMatrix(n, n)
	if n == 1 {
		(*ret)[0][0] = s1[0][0] * s2[0][0]
	} else {
		A11 := SplitMat(s1, n/2, n/2, 0, n/2, 0, n/2)
		A12 := SplitMat(s1, n/2, n/2, n/2, n, 0, n/2)
		A21 := SplitMat(s1, n/2, n/2, 0, n/2, n/2, n)
		A22 := SplitMat(s1, n/2, n/2, n/2, n, n/2, n)
		B11 := SplitMat(s2, n/2, n/2, 0, n/2, 0, n/2)
		B12 := SplitMat(s2, n/2, n/2, n/2, n, 0, n/2)
		B21 := SplitMat(s2, n/2, n/2, 0, n/2, n/2, n)
		B22 := SplitMat(s2, n/2, n/2, n/2, n, n/2, n)

		var M1, M2, M3, M4, M5, M6, M7 M
		doMultiplyStrassen(&M1, addMat(A11, A22), addMat(B11, B22))
		doMultiplyStrassen(&M2, addMat(A21, A22), B11)
		doMultiplyStrassen(&M3, A11, subMat(B12, B22))
		doMultiplyStrassen(&M4, A22, subMat(B21, B11))
		doMultiplyStrassen(&M5, addMat(A11, A12), B22)
		doMultiplyStrassen(&M6, subMat(A21, A11), addMat(B11, B12))
		doMultiplyStrassen(&M7, subMat(A12, A22), addMat(B21, B22))

		C11 := addMat(subMat(addMat(M1, M4), M5), M7)
		C12 := addMat(M3, M5)
		C21 := addMat(M2, M4)
		C22 := addMat(subMat(addMat(M1, M3), M2), M6)

		joinMat(*ret, C11, 0, n/2, 0, n/2)
		joinMat(*ret, C12, n/2, n, 0, n/2)
		joinMat(*ret, C21, 0, n/2, n/2, n)
		joinMat(*ret, C22, n/2, n, n/2, n)
	}
	return
}

func MultiplyStrassen(m1 M, m2 M) (ret M) {
	nr1, nc1 := m1.lens()
	nr2, nc2 := m2.lens()
	dim := calc2nDim(nr1, nc1, nr2, nc2)
	s1, s2 := MatToSqr(m1, m2, dim)
	doMultiplyStrassen(&ret, s1, s2)
	return SplitMat(ret, nr1, nc2, 0, nc2, 0, nr1)
}
