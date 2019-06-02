package main

import (
	"fmt"
	"log"
	"math"
	"utils"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getLossTrainning(tensor *[80][168][15]float64, U [][]float64, V [][]float64, W [][]float64, k int) (loss float64, lossTensor [80][168][15]float64) {
	loss = float64(0)
	lossTensor = [80][168][15]float64{}
	count := 0

	for si := 0; si < 80; si++ {
		for ti := 0; ti < 168; ti++ {
			for ci := 0; ci < 15; ci++ {
				if tensor[si][ti][ci] >= 0 {
					count++
					ev := float64(0)
					for ki := 0; ki < k; ki++ {
						ev += (U[si][ki]*V[ti][ki] + W[ci][ki]*V[ti][ki] + W[ci][ki]*U[si][ki])
					}
					delta := math.Abs(float64(ev - tensor[si][ti][ci]))

					loss += (delta * delta)
					lossTensor[si][ti][ci] = tensor[si][ti][ci] - ev
				} else {
					lossTensor[si][ti][ci] = -1.0
				}
			}
		}
	}

	return loss / float64(count), lossTensor
}

func getLossValidate(tensor [80][168][15]float64, tensorGT [80][168][15]float64, U [][]float64, V [][]float64, W [][]float64, k int) float64 {
	loss := float64(0)
	count := 0

	for si := 0; si < 80; si++ {
		for ti := 0; ti < 168; ti++ {
			for ci := 0; ci < 15; ci++ {
				if tensor[si][ti][ci] >= -1.1 && tensor[si][ti][ci] < 0 {
					count++
					ev := float64(0)
					for ki := 0; ki < k; ki++ {
						ev += (U[si][ki]*V[ti][ki] + W[ci][ki]*V[ti][ki] + W[ci][ki]*U[si][ki])
					}

					delta := math.Abs(float64(ev - tensorGT[si][ti][ci]))

					loss += (delta * delta)
				}
			}
		}
	}
	return loss / float64(count)
}

func gradientUpdateU(U *[][]float64, V [][]float64, W [][]float64, lossTensor *[80][168][15]float64, k int, tensor *[80][168][15]float64) {
	for si := 0; si < len(*U); si++ {
		for ki := 0; ki < k; ki++ {
			grad := float64(0)

			for ti := 0; ti < 168; ti++ {
				for ci := 0; ci < 15; ci++ {
					if tensor[si][ti][ci] >= 0 {
						grad += (lossTensor[si][ti][ci] * float64(V[ti][ki]+W[ci][ki]))
					}
				}
			}

			(*U)[si][ki] += float64(0.0000001 * grad)
		}
	}
}

func gradientUpdateV(V *[][]float64, U [][]float64, W [][]float64, lossTensor *[80][168][15]float64, k int, tensor *[80][168][15]float64) {
	for ti := 0; ti < len(*V); ti++ {
		for ki := 0; ki < k; ki++ {
			grad := float64(0)

			for si := 0; si < 80; si++ {
				for ci := 0; ci < 15; ci++ {
					if tensor[si][ti][ci] >= 0 {
						grad += (lossTensor[si][ti][ci] * float64(U[si][ki]+W[ci][ki]))
					}
				}
			}

			(*V)[ti][ki] += float64(0.00001 * grad)
		}
	}
}

func gradientUpdateW(W *[][]float64, U [][]float64, V [][]float64, lossTensor *[80][168][15]float64, k int, tensor *[80][168][15]float64) {
	for ci := 0; ci < len(*W); ci++ {
		for ki := 0; ki < k; ki++ {
			grad := float64(0)

			for si := 0; si < 80; si++ {
				for ti := 0; ti < 168; ti++ {
					if tensor[si][ti][ci] >= 0 {
						grad += (lossTensor[si][ti][ci] * float64(U[si][ki]+V[ti][ki]))
					}
				}
			}

			(*W)[ci][ki] += float64(0.000001 * grad)
		}
	}
}

func decompostion(tensor *[80][168][15]float64, k int) ([][]float64, [][]float64, [][]float64) {
	U := utils.InitalizeMatrix(80, k, true)
	V := utils.InitalizeMatrix(168, k, true)
	W := utils.InitalizeMatrix(15, k, true)

	loss, lossTensor := getLossTrainning(tensor, U, V, W, k)
	for loss > 0.1 {
		cacheU := utils.MatrixCopy(U)
		cacheV := utils.MatrixCopy(V)
		cacheW := utils.MatrixCopy(W)

		gradientUpdateU(&U, cacheV, cacheW, &lossTensor, k, tensor)
		gradientUpdateV(&V, cacheU, cacheW, &lossTensor, k, tensor)
		gradientUpdateW(&W, cacheU, cacheV, &lossTensor, k, tensor)

		loss, lossTensor = getLossTrainning(tensor, U, V, W, k)
		fmt.Println("loss", loss)
	}

	return U, V, W
}

func main() {
	fileNames := utils.FileNames()
	fmt.Println(len(fileNames))
	tensor := [80][168][15]float64{}
	tensorGT := [80][168][15]float64{}

	for _, fileName := range fileNames[1:] {
		fmt.Println(fileName)
		utils.ProcessFile(&tensor, &tensorGT, fileName)
	}
	utils.NormalizeByC(&tensor, &tensorGT)

	U, V, W := decompostion(&tensor, 30)
	loss := getLossValidate(tensor, tensorGT, U, V, W, 30)
	fmt.Println(loss)
}
