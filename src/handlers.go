package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
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

func (app *Application) UploadDocHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("Unable to parse multipart form: %v\n", err)
		return
	}

	if len(r.MultipartForm.File["file"]) != 1 {
		log.Print("Invalid File Input")
		return
	}

	hdr := r.MultipartForm.File["file"][0]
	f, err := hdr.Open()
	if err != nil {
		log.Printf("Unable to open file header's assoicated file: %v\n", err)
		return
	}
	defer f.Close()

	ext := filepath.Ext(hdr.Filename)
	var id int

	id = rand.Intn(100)

	oldname := fmt.Sprintf("uploads/%s_%s", strconv.Itoa(id), hdr.Filename)

	file, err := os.Create(oldname)
	if err != nil {
		file.Close()
		log.Printf("Unable to create file in uploads err: %v\n", err)
		return
	}

	io.Copy(file, f)

	file.Close()

	nff, err := os.Open(oldname)
	if err != nil {
		log.Printf("Unable to open file in uploads err: %v\n", err)
		return
	}
	defer nff.Close()

	switch ext {
	case ".pdf", ".docx":
		id, err = app.InsertDoc(nff, hdr.Filename)
	case ".txt":
		id = app.Insert(nff, hdr.Filename)
	}

	if err != nil {
		log.Printf("Unable to insert file to database: %v\n", err)
		return
	}

	if err := os.Rename(oldname, fmt.Sprintf("uploads/%s%s", strconv.Itoa(id), ext)); err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	// http.Redirect(w, r, "/", http.StatusAccepted)
}
