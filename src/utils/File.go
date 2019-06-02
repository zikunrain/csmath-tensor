package utils

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FileNames() []string {
	var files []string

	root := "src/data"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	check(err)
	return files
}

func generateMap() map[string]int {
	var m map[string]int /*创建集合 */
	m = make(map[string]int)
	m["AQI"] = 0
	m["PM2.5"] = 1
	m["PM2.5_24h"] = 2
	m["PM10"] = 3
	m["PM10_24h"] = 4
	m["SO2"] = 5
	m["SO2_24h"] = 6
	m["NO2"] = 7
	m["NO2_24h"] = 8
	m["O3"] = 9
	m["O3_24h"] = 10
	m["O3_8h"] = 11
	m["O3_8h_24h"] = 12
	m["CO"] = 13
	m["CO_24h"] = 14
	return m
}

func NormalizeByC(tensor *[80][168][15]float64, tensorGT *[80][168][15]float64) {
	cMax := []float64{}
	for ci := 0; ci < 15; ci++ {
		max := float64(0.0)
		for si := 0; si < 80; si++ {
			for ti := 0; ti < 168; ti++ {
				if tensor[si][ti][ci] > max {
					max = tensor[si][ti][ci]
				}
			}
		}
		cMax = append(cMax, max)
	}
	for si := 0; si < 80; si++ {
		for ti := 0; ti < 168; ti++ {
			for ci := 0; ci < 15; ci++ {
				if tensorGT[si][ti][ci] >= 0 {
					tensorGT[si][ti][ci] = tensorGT[si][ti][ci] / cMax[ci]

					if rand.Float64() > 0.8 {
						tensor[si][ti][ci] = -1.0
					} else {
						tensor[si][ti][ci] = tensor[si][ti][ci] / cMax[ci]
					}
				}
			}
		}
	}
}

func ProcessFile(tensor *[80][168][15]float64, tensorGT *[80][168][15]float64, fileName string) {
	cMap := generateMap()

	file, err := os.Open(fileName)
	check(err)

	fmt.Println(fileName)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	l := 0
	for scanner.Scan() {
		l++
		if l == 1 {
			continue
		}

		params := strings.Split(scanner.Text(), ",")
		t, err := strconv.Atoi(params[1])
		date, err := strconv.Atoi(params[0])
		t += (date - 20141224) * 24
		check(err)
		c, _ := cMap[params[2]]
		values := params[3:]
		for si, v := range values {
			if si >= 80 {
				break
			}
			fv, err := strconv.ParseFloat(v, 32)
			if err != nil {
				tensor[si][t][c] = float64(-2)
				tensorGT[si][t][c] = float64(-2)
			} else {
				tensor[si][t][c] = float64(fv)
				tensorGT[si][t][c] = float64(fv)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
