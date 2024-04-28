package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"gonum.org/v1/gonum/floats"
)

// Map of each word in the embeddings file to it's corresponding vector
var WordVec map[string][]float64 = make(map[string][]float64)

// Initialize the word vector map from embeddings file
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

// Sum of embedding vector for each line in the document
func CalcDocVec(in io.Reader) []float64 {
	res := make([]float64, 10) // an array of 10 floats, 10 is the dimension of the embeddings
	var vec []float64

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()

		if line == " " {
			continue
		}

		vec = CalcQueryVec(line)
		floats.AddTo(res, res, vec)
	}

	return res
}

func NCalcDocVec(buf *bytes.Buffer) []float64 {
	rd := bufio.NewReader(buf)
	rs := make([]float64, 10)
	for {
		wd, err := rd.ReadString(byte(' '))
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		op := strings.ToLower(sanitizeWord(wd))
		vec, ok := WordVec[op]

		if !ok {
			continue
		}
		floats.AddTo(rs, rs, vec)
	}

	return rs
}

// Sum of embedding vectors for each word in the line
func CalcQueryVec(query string) []float64 {
	res := make([]float64, 10)
	words := strings.Split(query, " ")
	for _, word := range words {
		op := strings.ToLower(sanitizeWord(word))
		vec, ok := WordVec[op]
		if !ok {
			continue
		}
		floats.AddTo(res, res, vec)
	}
	return res
}

// Calculates the approximate similiarity of two vectors
func CosineSim(a, b []float64) float64 {
	na := floats.Norm(a, 2)
	nb := floats.Norm(b, 2)

	return 1 - (floats.Dot(a, b) / (na * nb))
}

func GetVec(word string) ([]float64, bool) {
	x, ok := WordVec[word]
	return x, ok
}
