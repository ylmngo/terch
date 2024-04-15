package application

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strconv"
	"terch/utils"

	"github.com/google/go-tika/tika"
	"github.com/lib/pq"
)

type Application struct {
	Db      *sql.DB
	TikaCli *tika.Client
}

type Document struct {
	Id   int
	Name string
	Vec  []float64
}

type DocumentResult struct {
	DocId int
	Sim   float64
}

var App Application

func InitApp(file *os.File, dsn string) *Application {
	utils.InitMap(file)
	App.connectDB(dsn)
	App.TikaCli = utils.InitTikaClient("http://localhost:9998")

	return &App
}

func (app *Application) Search(query string) ([]DocumentResult, error) {
	res := make([]DocumentResult, 0)

	docs, err := app.GetAllFromDB()
	if err != nil {
		return nil, err
	}

	a := utils.CalcQueryVec(query)
	for _, d := range docs {
		fmt.Println(d.Vec)
		sim := utils.CosineSim(d.Vec, a)
		res = append(res, DocumentResult{DocId: d.Id, Sim: sim})
	}

	slices.SortFunc(res, func(a, b DocumentResult) int {
		if a.Sim < b.Sim {
			return -1
		} else if a.Sim > b.Sim {
			return 1
		}
		return 0
	})

	l := len(res)
	if len(res) > 5 {
		l = 5
	}
	return res[:l], nil
}

func (app *Application) Insert(f *os.File) error {
	var id int

	vec := utils.CalcDocVec(f)

	app.Db.QueryRow(`INSERT into docs (name, docvec) VALUES ($1, $2) RETURNING id`, f.Name(), pq.Array(vec)).Scan(&id)

	file, err := os.Create(fmt.Sprintf("uploads/%s.txt", strconv.Itoa(id)))
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

// Inserts Name of the file and it's document vector to database and saves it to disk
func (app *Application) InsertPDF(f *os.File) (int, error) {
	var id int

	doc, err := utils.CreateDocument(f, app.TikaCli)
	if err != nil {
		return 0, err
	}

	app.Db.QueryRow(`INSERT into docs (name, docvec) VALUES ($1, $2) RETURNING id`, doc.Name, pq.Array(doc.Vec)).Scan(&id)

	file, err := os.Create(fmt.Sprintf("uploads/%s.txt", strconv.Itoa(id)))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return id, nil
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

func (app *Application) GetAllFromDB() ([]Document, error) {
	res := make([]Document, 0)
	rows, err := app.Db.Query(`SELECT * FROM docs`)
	if err != nil {
		return nil, fmt.Errorf("unable to get rows from database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		d := Document{}
		if err := rows.Scan(&d.Id, &d.Name, pq.Array(&d.Vec)); err != nil {
			fmt.Printf("Unable to scan rows to document struct: %v\n", err)
			continue
		}
		res = append(res, d)
	}

	return res, nil
}

func (app *Application) connectDB(DSN string) {
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Unable to Ping to databse: %v\n", err)
		return
	}

	app.Db = db
}
