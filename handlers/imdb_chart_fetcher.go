package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type ImdbChartFetcher struct {
	l *log.Logger
}

func NewImdbChartFetcher(l *log.Logger) *ImdbChartFetcher {
	return &ImdbChartFetcher{l}

}
func (i *ImdbChartFetcher) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle the request for a list of products
	if r.Method == http.MethodGet {
		i.fetchImdbChart(rw, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	rw.WriteHeader(http.StatusMethodNotAllowed)
}
func (i *ImdbChartFetcher) getHref(t html.Token) (ok bool, href string) {
	// Iterate over token attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" && strings.Contains(a.Val, "title") {
			href = a.Val
			ok = true
		}
	}

	return
}
func (i *ImdbChartFetcher) fetchImdbChart(rw http.ResponseWriter, r *http.Request) {
	i.l.Println("Handle GET Products")

	url := "https://www.imdb.com/india/top-rated-indian-movies/"
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	//scanner := bufio.NewScanner(resp.Body)
	b := resp.Body
	defer b.Close() // close Body when the function completes

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()
		//	b := z.Text()
		//fmt.Println("ttt----", tt.Data)
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		}
		t := z.Token()

		// Check if the token is an <a> tag

		isAnchor := t.Data == "a"

		if !isAnchor {
			continue
		}

		// Extract the href value, if there is one
		ok, url := i.getHref(t)
		if !ok {

			continue
		}

		// Make sure the url begines in http**
		hasProto := strings.Index(url, "http") == 0

		fmt.Println("link------", url, hasProto)

	}
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
