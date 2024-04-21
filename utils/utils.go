package utils

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"
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

// Try removing regexp
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

func CosineSim(a, b []float64) float64 {
	na := floats.Norm(a, 2)
	nb := floats.Norm(b, 2)

	return 1 - (floats.Dot(a, b) / (na * nb))
}

func GetVec(word string) ([]float64, bool) {
	x, ok := WordVec[word]
	return x, ok
}

func sanitizeWord(word string) string {
	var builder strings.Builder
	for _, s := range word {
		if !unicode.IsLetter(s) {
			continue
		}
		builder.WriteRune(s)
	}
	return builder.String()
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

func ParseTemplate(path string) *template.Template {
	tmpl, _ := template.ParseFiles(path)
	return tmpl
}

func ConnectDBPool(DSN string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), DSN)
	if err != nil {
		return pool, nil
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil

	// db, err := sql.Open("postgres", DSN)
	// if err != nil {
	// 	return nil, err
	// }

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// if err := db.PingContext(ctx); err != nil {
	// 	return nil, err
	// }

	// return db, nil
}
