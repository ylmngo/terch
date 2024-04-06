package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var WordVector map[string][]float64 = make(map[string][]float64, 1)

const BUFFER_SIZE int = 2048 * 2048

func main() {
	file, err := os.Open("oup_prc.txt")
	if err != nil {
		fmt.Printf("Unable to open word vector file: %v\n", err)
		return
	}
	defer file.Close()
	initMap(file)
}

func initMap(f *os.File) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		word, vec, err := parseLine(line)
		if err != nil {
			continue
		}
		WordVector[word] = vec
	}
}

func parseLine(line string) (string, []float64, error) {
	parts := strings.Split(line, " ")
	word := parts[0]
	vec := make([]float64, 10)
	for i := 1; i < 11; i++ {
		val, err := strconv.ParseFloat(parts[i], 64)
		if err != nil {
			fmt.Printf("Unable to parse floating point value: %v\n", err)
			return "", nil, err
		}
		vec[i-1] = val
	}
	return word, vec, nil
}
