package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"terch/utils"
)

// Pre parse the templates and then execute them

var (
	searchTempl *template.Template = utils.ParseTemplate("templates/search.html")
	indexTempl  *template.Template = utils.ParseTemplate("templates/index.html")
)

func (app *Application) SearchHandler(w http.ResponseWriter, r *http.Request) {
	vals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Printf("Unable to parse query parameters from url: %s\n, err: %v\n", r.URL.Path, err)
		return
	}

	query := vals.Get("query")
	if query == "" {
		log.Printf("No query parameters in search path: %s\n", query)
		return
	}

	res, err := app.Search(query)
	if err != nil {
		log.Printf("Unable retrieve documents for query: %s\n, %v\n", query, err)
		fmt.Fprintf(w, "Cannot find results for query: %s\n", query)
		return
	}

	if err := searchTempl.Execute(w, res); err != nil {
		log.Printf("Unable to execute template file: %s, err: %v\n", "templates/search.html", err)
		return
	}

}

func (app *Application) IndexHandler(w http.ResponseWriter, r *http.Request) {

	if err := indexTempl.Execute(w, nil); err != nil {
		log.Printf("Unable to execute template file: %s, err: %v\n", "templates/index.html", err)
		return
	}
}

func (app *Application) ViewDocHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Printf("Invalid Document ID: %d\n", id)
		return
	}

	filename, err := app.GetDocument(id)
	if err != nil {
		log.Printf("Unable to get filename for id: %d\n, err: %v\n", id, err)
		return
	}

	ext := filepath.Ext(filename)

	file, err := os.Open(fmt.Sprintf("uploads/%d%s", id, ext))
	if err != nil {
		log.Printf("Unable to open file: uploads/%d, err: %v\n", id, err)
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	if _, err = file.Read(buf); err != nil {
		log.Printf("Unable to read to buf from file: uploads/%d\n, err: %v\n", id, err)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Unable to read contents of file: uploads/%d, err: %v\n", id, err)
		return
	}

	ctype := http.DetectContentType(buf)
	w.Header().Add("Content-Type", ctype)
	w.Write(data)
}

func (app *Application) UploadFileHandler(w http.ResponseWriter, r *http.Request) {

}
