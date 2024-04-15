package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/floats"
)

var WordVec map[string][]float64 = make(map[string][]float64)

func InitMap(f *os.File) {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		word, vec, err := parseLine(line)
		if err != nil {
			continue
		}
		WordVec[word] = vec
	}
}

func CalcDocVec(f *os.File) []float64 {
	res := make([]float64, 10) // an array of 10 floats, 10 is the dimension of the embeddings

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if line == " " {
			continue
		}

		vec := CalcQueryVec(line)
		floats.AddTo(res, res, vec)
	}

	return res
}

func NCalcDocVec(buf *bytes.Buffer) []float64 {
	rd := bufio.NewReader(buf)
	re := regexp.MustCompile(`[^\d\p{Latin}]`)
	rs := make([]float64, 10)
	for {
		wd, err := rd.ReadString(byte(' '))
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		op := strings.ToLower(re.ReplaceAllString(wd, ""))
		vec, ok := WordVec[op]

		if !ok {
			continue
		}
		floats.AddTo(rs, rs, vec)
	}

	return rs
}

// Try removing regexp
func CalcQueryVec(query string) []float64 {
	res := make([]float64, 10)
	re := regexp.MustCompile(`[^\d\p{Latin}]`)
	words := strings.Split(query, " ")
	for _, word := range words {
		op := strings.ToLower(re.ReplaceAllString(word, ""))
		vec, ok := WordVec[op]
		if !ok {
			continue
		}
		floats.AddTo(res, res, vec)
	}

	return res
}

func CosineSim(a, b []float64) float64 {
	na := floats.Norm(a, 2)
	nb := floats.Norm(b, 2)

	return 1 - (floats.Dot(a, b) / (na * nb))
}

func GetVec(word string) ([]float64, bool) {
	x, ok := WordVec[word]
	return x, ok
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
