package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"terch/utils"

	_ "net/http/pprof"

	_ "github.com/lib/pq"

	"github.com/google/go-tika/tika"
)

type Application struct {
	Db      *sql.DB
	TikaCli *tika.Client
	Router  *http.ServeMux
}

const DSN string = "postgres://terch:freeroam@localhost/terchdb?sslmode=disable"

func main() {
	file, err := os.Open("oup.txt")
	if err != nil {
		log.Fatalf("Unable to find embeddings file: %v\n", err)
	}
	defer file.Close()

	app := InitApp(file, DSN)
	app.Router = app.Routes()

	srv := &http.Server{
		Addr:    ":8000",
		Handler: app.Router,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server Failure: %v\n", err)
	}
}

func InitApp(file *os.File, dsn string) *Application {
	app := &Application{}

	utils.InitMap(file)
	app.ConnectDB(dsn)
	app.TikaCli = utils.InitTikaClient("http://localhost:9998")

	return app
}

func (app *Application) ConnectDB(DSN string) {
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
