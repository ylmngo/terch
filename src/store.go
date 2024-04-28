package main

import "time"

// Allocates a new map for storing all of the contents of the database
// Locks the store during it's execution
func (app *Application) StoreDocuments() {
	app.Store.Lock()
	defer app.Store.Unlock()
	app.Store.data = make(map[int]Document)
	docs, _ := app.GetAllFromDB()
	for _, doc := range docs {
		app.Store.data[doc.Id] = doc
	}
}

// Get all the documents from store
func (app *Application) GetAlls() []Document {
	docs := make([]Document, len(app.Store.data))
	i := 0
	for id, doc := range app.Store.data {
		docs[i] = Document{Id: id, Name: doc.Name, Vec: doc.Vec}
		i += 1
	}
	return docs
}

// Re-allocates the store every 20 seconds
func (app *Application) RestoreDocuments() {
	for {
		time.Sleep(time.Until(<-time.After(20 * time.Second)))
		app.StoreDocuments()
	}
}
