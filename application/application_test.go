package application_test

import (
	"os"
	"terch/application"
	"testing"
)

func TestInsertDocument(t *testing.T) {
	file, _ := os.Open("oup_prc.txt")
	defer file.Close()

	app := application.InitApp(file, "postgres://terch:freeroam@localhost/terchdb?sslmode=disable")

	pd, _ := os.Open("bgnet.pdf")
	defer pd.Close()

	id, err := app.InsertPDF(pd)
	if err != nil {
		t.Fatalf("Unable to insert PDF file to Database: %v\n", err)
		return
	}

	if _, err := app.Db.Exec(`DELETE FROM docs WHERE id = $1`, id); err != nil {
		t.Fatalf("unable to delete row from database: %v\n", err)
		return
	}
}
