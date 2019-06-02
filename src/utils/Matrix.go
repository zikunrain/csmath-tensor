package utils

import "math/rand"

func MatrixCopy(m [][]float64) [][]float64 {
	d := InitalizeMatrix(len(m), len(m[0]), false)
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[0]); j++ {
			d[i][j] = m[i][j]
		}
	}
	return d
}

func InitalizeMatrix(x int, k int, isRand bool) [][]float64 {
	matrix := make([][]float64, x)
	for i := 0; i < x; i++ {
		matrix[i] = make([]float64, k)
	}
	if isRand {
		for i := 0; i < x; i++ {
			for j := 0; j < k; j++ {
				matrix[i][j] = rand.Float64()
			}
		}
		return matrix
	} else {
		return matrix
	}
}
