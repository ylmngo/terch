package main

// import (
// 	"bufio"
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"os"
// 	"regexp"
// 	"strconv"
// 	"strings"

// 	"github.com/lib/pq"
// 	"gonum.org/v1/gonum/floats"
// )

// var WordVector map[string][]float64 = make(map[string][]float64, 1)

// type DocSim struct {
// 	name string
// 	val  float64
// }

// const BUFFER_SIZE int = 2048 * 2048
// const DSN string = "postgres://terch:freeroam@localhost/terchdb?sslmode=disable"

// func main() {
// 	file, err := os.Open("oup.txt")
// 	if err != nil {
// 		fmt.Printf("Unable to open word vector file: %v\n", err)
// 		return
// 	}
// 	defer file.Close()

// 	initMap(file)

// 	db, err := openDB(DSN)
// 	if err != nil {
// 		fmt.Printf("unable to create connection to the database: %v\n", err)
// 		return
// 	}
// 	defer db.Close()

// 	txtFile, err := os.Open("uploads/history.txt")
// 	if err != nil {
// 		fmt.Printf("unable to open uploads file: %v\n", err)
// 		return
// 	}
// 	defer txtFile.Close()

// 	addToDB(txtFile, db)

// 	fmt.Println(calcQueryDocSims("Firewalls, encryption and user detection software algorithms", db))

// 	// a := calcDocVec(txtFile)
// 	// b := calcQueryVec("crucial to protect data from unauthorized access and attacks")
// 	// na := floats.Norm(a, 2)
// 	// nb := floats.Norm(b, 2)
// 	// cosSim := 1 - floats.Dot(a, b)/(na*nb)
// 	// fmt.Println(cosSim)
// }

// func calcQueryDocSims(query string, db *sql.DB) []DocSim {
// 	res := make([]DocSim, 0, 5)
// 	rows, err := db.Query(`SELECT * FROM docs`)
// 	if err != nil {
// 		fmt.Printf("Unable to Retrieve values from DB: %v\n", err)
// 		return nil
// 	}
// 	defer rows.Close()

// 	var id int
// 	var name string

// 	for rows.Next() {
// 		var vecBuf []float64 = make([]float64, 1)
// 		if err := rows.Scan(&id, &name, pq.Array(&vecBuf)); err != nil {
// 			fmt.Printf("Unable to scan row: %v\n", err)
// 			continue
// 		}

// 		d := DocSim{name: name, val: cosineSim(query, vecBuf)}
// 		res = append(res, d)
// 	}

// 	return res
// }

// func cosineSim(query string, vec []float64) float64 {
// 	a := calcQueryVec(query)
// 	b := vec

// 	na := floats.Norm(a, 2)
// 	nb := floats.Norm(b, 2)

// 	return (1 - floats.Dot(a, b)/(na*nb))
// }

// func initMap(f *os.File) {
// 	scanner := bufio.NewScanner(f)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		word, vec, err := parseLine(line)
// 		if err != nil {
// 			continue
// 		}
// 		WordVector[word] = vec
// 	}
// }

// func parseLine(line string) (string, []float64, error) {
// 	parts := strings.Split(line, " ")
// 	word := parts[0]
// 	vec := make([]float64, 10)
// 	for i := 1; i < 11; i++ {
// 		val, err := strconv.ParseFloat(parts[i], 64)
// 		if err != nil {
// 			fmt.Printf("Unable to parse floating point value: %v\n", err)
// 			return "", nil, err
// 		}
// 		vec[i-1] = val
// 	}
// 	return word, vec, nil
// }

// func calcQueryVec(query string) []float64 {
// 	re := regexp.MustCompile(`[^\d\p{Latin}]`)
// 	words := strings.Split(query, " ")
// 	dst := make([]float64, 10)
// 	for _, word := range words {
// 		res := re.ReplaceAllString(word, "")
// 		res = strings.ToLower(res)
// 		vec, ok := WordVector[res]
// 		// Lemmetization of words

// 		if !ok {
// 			fmt.Println(res)
// 			continue
// 		}
// 		floats.AddTo(dst, dst, vec)
// 	}
// 	return dst
// }

// func openDB(dsn string) (*sql.DB, error) {
// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	if err = db.PingContext(ctx); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

// // TODO: save file to disk
// func addToDB(f *os.File, db *sql.DB) {
// 	vec := calcDocVec(f)
// 	name := f.Name()
// 	_, err := db.Exec(`INSERT into docs (name, docvec) VALUES ($1, $2)`, name, pq.Array(vec))
// 	if err != nil {
// 		fmt.Printf("Error while executing DB Insertion: %v\n", err)
// 		return
// 	}
// 	defer f.Close()
// }

// func calcDocVec(f *os.File) []float64 {
// 	scanner := bufio.NewScanner(f)
// 	res := make([]float64, 10)
// 	for scanner.Scan() {
// 		line := scanner.Text()

// 		if line == " " {
// 			continue
// 		}

// 		vec := calcQueryVec(line)

// 		floats.AddTo(res, res, vec)
// 	}

// 	return res
// }
