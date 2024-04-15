package main

import (
	"net/http"
	"net/http/pprof"
)

func (app *Application) Routes() *http.ServeMux {
	router := &http.ServeMux{}

	router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	router.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	router.HandleFunc("/index", app.IndexHandler)
	router.HandleFunc("/search", app.SearchHandler)
	router.HandleFunc("/document/{id}", app.ViewDocHandler)
	// router.HandleFunc("/upload", s.UploadDocHandler)
	// router.HandleFunc("/delete", s.DeleteDocHandler)

	return router
}
