package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"terch/utils"

	_ "net/http/pprof"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"github.com/google/go-tika/tika"
)

type Application struct {
	DbPool  *pgxpool.Pool
	TikaCli *tika.Client
	Router  *http.ServeMux
	Store   struct {
		sync.Mutex
		data map[int]Document
	}
}

const DSN string = "postgres://terch:freeroam@localhost/terchdb?sslmode=disable"

func main() {
	file, err := os.Open("oup.txt")
	if err != nil {
		log.Fatalf("Unable to find embeddings file: %v\n", err)
	}
	defer file.Close()

	app := InitApp(file, DSN)
	defer app.DbPool.Close()

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
	app.DbPool, _ = utils.ConnectDBPool(dsn)
	app.Router = app.Routes()
	app.TikaCli = utils.InitTikaClient("http://localhost:9998")

	return app
}
