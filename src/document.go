package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"terch/utils"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/lib/pq"
)

type Document struct {
	Id   int
	Name string
	Vec  []float64
}

type DocumentResult struct {
	Name  string
	DocId int
	Sim   float64
}

// Change it such that it returns an id, the calling function shall be responsible for creating upload file
func (app *Application) Insert(f *os.File) error {
	var id int

	vec := utils.CalcDocVec(f)

	app.DbPool.QueryRow(context.Background(), `INSERT into docs (name, docvec) VALUES ($1, $2) RETURNING id`, f.Name(), pq.Array(&vec)).Scan(&id)

	file, err := os.Create(fmt.Sprintf("uploads/%s.txt", strconv.Itoa(id)))
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// Inserts Name of the file and it's document vector to database and saves it to disk
// Change it such that it returns an id, the calling function shall be responsible for creating upload file

func (app *Application) InsertPDF(f *os.File) (int, error) {
	var id int

	doc, err := utils.CreateDocument(f, app.TikaCli)
	if err != nil {
		return 0, err
	}

	app.DbPool.QueryRow(context.Background(), `INSERT into docs (name, docvec) VALUES ($1, $2) RETURNING id`, doc.Name, pq.Array(doc.Vec)).Scan(&id)

	file, err := os.Create(fmt.Sprintf("uploads/%s.txt", strconv.Itoa(id)))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return id, nil
}

func (app *Application) GetAllFromDB() ([]Document, error) {
	rows, err := app.DbPool.Query(context.Background(), `SELECT id, name, docvec FROM docs`)
	if err != nil {
		return nil, fmt.Errorf("unable to get rows from database: %v", err)
	}
	defer rows.Close()

	res := make([]Document, 0)

	d := Document{}
	for rows.Next() {
		if err := rows.Scan(&d.Id, &d.Name, &d.Vec); err != nil {
			fmt.Printf("Unable to scan rows to document struct: %v\n", err)
			continue
		}
		res = append(res, d)
	}

	return res, nil
}

func (app *Application) GetDocument(id int) (string, error) {
	row := app.DbPool.QueryRow(context.Background(), `SELECT name FROM docs WHERE id = $1`, id)
	var name string
	if err := row.Scan(&name); err != nil {
		return "", err
	}
	return name, nil
}

// func (app *Application) InsertDocument(doc *utils.Document) error {
// 	var id int

// 	app.Db.QueryRow(`INSERT into docs (name, docvec) VALUES ($1, $2) RETURNING id`, doc.Name, doc.Vec).Scan(&id)

// 	file, err := os.Create(fmt.Sprintf("uploads/%s.txt", strconv.Itoa(id)))
// 	if err != nil {
// 		return errors.New("unable to create file in uploads dir")
// 	}
// 	defer file.Close()

// 	return nil
// }
