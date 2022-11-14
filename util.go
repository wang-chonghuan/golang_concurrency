package main

func max(a int, b int) int {
	if a > b {
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
