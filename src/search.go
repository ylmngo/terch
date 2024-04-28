package main

import (
	"slices"
	"terch/utils"
)

func (app *Application) Search(query string) ([]DocumentResult, error) {
	docs := app.GetAlls()

	res := make([]DocumentResult, 0, len(docs))

	a := utils.CalcQueryVec(query)
	for _, d := range docs {
		sim := utils.CosineSim(d.Vec, a)
		res = append(res, DocumentResult{Name: d.Name, DocId: d.Id, Sim: sim})
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
