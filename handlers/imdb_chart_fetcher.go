package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/wajidzaman/imdb_fetcher/data"
)

const (
	IMDB_PREFIX    = "https://www.imdb.com"
	IMDB_URL_REGEX = `<a href=\"/title/tt(.*?)"`
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

func (i *ImdbChartFetcher) fetchImdbChart(rw http.ResponseWriter, r *http.Request) {
	i.l.Println("Handle GET Products")
	url := "https://www.imdb.com/india/top-rated-indian-movies/"
	urls := i.getMovieUrl(url)
	i.parseMovieFromUrlsConcurrently(urls, 1)

}

func (i *ImdbChartFetcher) parseMovieFromUrlsConcurrently(urls []string, k int) {
	if len(urls) < k {
		k = len(urls)
	}
	movie := make(chan data.Movie, k)
	var wg sync.WaitGroup
	for j := 0; j < k; j++ {
		go i.parseEachUrl(&wg, j, IMDB_PREFIX+urls[j], movie)
	}
	wg.Wait()
	var movies []*data.Movie
	for m := range movie {
		movies = append(movies, &m)
	}
}

func (i *ImdbChartFetcher) parseEachUrl(wg *sync.WaitGroup, movieSNo int, url string, movieChan chan data.Movie) {

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	bodyHtml, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("-----------")
	fmt.Println("")
	movie := data.GetMoviesByParsingHTML(string(bodyHtml))
	movieChan <- *movie
	fmt.Println("movies:", movie)
	wg.Done()
}

func (i *ImdbChartFetcher) getMovieUrl(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	b, _ := ioutil.ReadAll(resp.Body)
	body := string(b)

	// parsing url from html for getting movies details
	body = strings.Replace(body, "\n", "", -1)
	re := regexp.MustCompile(IMDB_URL_REGEX)
	rawUrls := re.FindAllString(body, -1)
	m := make(map[string]bool)
	var urls []string
	for _, e := range rawUrls {
		if m[e] == false {
			r, _ := regexp.Compile("=\"")
			index := r.FindStringIndex(e)
			url := e[index[0]+2 : len(e)-1]
			m[e] = true
			urls = append(urls, url)

		}

	}
	for _, url := range urls {
		fmt.Println("url:", url)
	}

	return urls
}
