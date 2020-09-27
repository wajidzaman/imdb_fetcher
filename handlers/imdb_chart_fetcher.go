package handlers

import (
	"log"
	"net/http"
)

type ImdbChartFetcher struct {
	l *log.Logger
}

func NewImdbChartFetcher(l *log.Logger) *ImdbChartFetcher {
	return &ImdbChartFetcher{l}
}
func (i *ImdbChartFetcher) fetchImdbChart(rw http.ResponseWriter, r *http.Request) {
	i.l.Println("Handle GET Products")

	// fetch the products from the datastore
	/*
		lp := data.GetProducts()

		// serialize the list to JSON
		err := lp.ToJSON(rw)
		if err != nil {
			http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
		}
	*/
}
