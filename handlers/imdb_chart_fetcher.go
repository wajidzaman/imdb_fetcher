package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
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
	inputUrl, ok := r.URL.Query()["url"]

	if !ok || len(inputUrl[0]) < 1 {
		log.Println("url is missing")
		return
	}
	url := inputUrl[0]
	movieCount, ok := r.URL.Query()["k"]
	k, err := strconv.Atoi(movieCount[0])
	if !ok || len(movieCount[0]) < 1 || err != nil {
		log.Println("k is missing")
		return
	}

	urls := i.getMovieUrl(url)

	i.parseMovieFromUrlsConcurrently(rw, urls, k)

}

func (i *ImdbChartFetcher) parseMovieFromUrlsConcurrently(rw http.ResponseWriter, urls []string, k int) error {
	if len(urls) < k {
		k = len(urls)
	}
	movie := make(chan data.Movie, k)
	jobOrder := make(chan int, k)

	var wg sync.WaitGroup
	for j := 0; j < k; j++ {
		wg.Add(1) //maintaining the number of goroutines we are firing
		go i.parseEachUrlJob(&wg, j, IMDB_PREFIX+urls[j], movie, jobOrder)
	}
	wg.Wait() // waiting for all gorutine to complete there processing.
	movieOrderMap := make(map[int]data.Movie)

	var movies []*data.Movie
	for j := 0; j < k; j++ {
		movideDetail := <-movie
		movieOrderId := <-jobOrder
		movieOrderMap[movieOrderId] = movideDetail
	}
	for j := 0; j < k; j++ {
		movie := movieOrderMap[j]
		movies = append(movies, &movie)
	}
	i.l.Println("movies return", len(movies))
	e := json.NewEncoder(rw)
	return e.Encode(movies)
}

func (i *ImdbChartFetcher) parseEachUrlJob(wg *sync.WaitGroup, movieSNo int, url string, movieChan chan data.Movie, jobOrder chan int) {
	i.l.Println("fetching url : ", url)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	i.l.Println("Url :", url, " Response status:", resp.Status)

	bodyHtml, _ := ioutil.ReadAll(resp.Body)

	movie := data.GetMoviesByParsingHTML(string(bodyHtml))
	movieChan <- movie
	jobOrder <- movieSNo
	//fmt.Println("movies:", movie)

	wg.Done()
}

func (i *ImdbChartFetcher) getMovieUrl(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	i.l.Println("Response status:", resp.Status)

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
	return urls
}
