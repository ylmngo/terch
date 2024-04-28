package main

import "time"

func (app *Application) StoreDocuments() {
	app.Store.Lock()
	defer app.Store.Unlock()
	app.Store.data = make(map[int]Document)
	docs, _ := app.GetAllFromDB()
	for _, doc := range docs {
		app.Store.data[doc.Id] = doc
	}
}

func (app *Application) GetAlls() []Document {
	docs := make([]Document, len(app.Store.data))
	i := 0
	for id, doc := range app.Store.data {
		docs[i] = Document{Id: id, Name: doc.Name, Vec: doc.Vec}
	}
	return docs
}

func (app *Application) ReStoreDocuments() {
	for {
		time.Sleep(time.Until(<-time.After(20 * time.Second)))
		app.StoreDocuments()
	}
}
