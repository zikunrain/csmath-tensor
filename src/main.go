package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// type Tensor struct {
// 	// LenS   int
// 	// LenT   int
// 	// LenC   int
// 	Values [80][168][14]float32
// }

// func NewTensor(lenS int, lenT int, lenC int) *Tensor {
// 	t := new(Tensor)
// 	t.LenS = lenS
// 	t.LenT = lenT
// 	t.LenC = lenC // <- a very sensible default value
// 	t.Values =
// 	return t
// }

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

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func FileNames() []string {
	var files []string

	root := "data"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	check(err)
	return files
}

func ProcessFile(tensor [80][168][14]float32, fileName string, cMap map[string]int) {
	file, err := os.Open(fileName)
	check(err)

	fmt.Println(fileName)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	l := 0
	for scanner.Scan() {
		l++
		if l == 1 { // header
			continue
		}

		params := strings.Split(scanner.Text(), ",")
		t, err := strconv.Atoi(params[1])
		check(err)
		c, _ := cMap[params[2]]
		values := params[3:]
		for si, v := range values {
			fv, err := strconv.ParseFloat(v, 32)
			check(err) // if err set -1?
			tensor[si][t][c] = float32(fv)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fileNames := FileNames()
	// s := 80
	// t := 24 * 7
	// c := 14
	tensor := [80][168][14]float32{}
	cMap := generateMap()

	for j, fileName := range fileNames[1:] {
		if j%100 == 0 {
			fmt.Println(j, len(fileNames))
		}
		ProcessFile(tensor, fileName, cMap)
		// process current file
	}
}
