package utils

import (
	"context"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
}

// Removes invalid words - alphanumeric, Non-English characters...
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
