package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

const (
	IMDB_GENRE    = "\"itemprop\" itemprop=\"genre\">(.*?)</span>"
	IMDB_SUMMARY  = `<div class=\"summary_text\">(.*?)<\/div>`
	IMDB_RATING   = "<span itemprop=\"ratingValue\">(.*?)</span>"
	IMDB_TITLE    = "property='og:title' content=\"(.*?)>"
	IMDB_YEAR     = "property='og:title' content=\".*?\\(([0-9]{4})\\)\""
	IMDB_GENRES   = `<h4 class=\"inline\">Genres:</h4>(.*?)<\/div>`
	IMDB_DURATION = " <h4 class=\"inline\">Runtime:</h4>(.*?)</time>"
)

type Movie struct {
	Title    string
	Year     int64
	Rating   float64
	Summary  string
	Duration string
	Genre    string
}

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
type Movies []*Movie

func (m *Movies) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

// GetProducts returns a list of products

func GetMoviesByParsingHTML(body string) *Movie {
	var title string
	var year int64
	var rating float64
	var summary string
	var duration string
	var genreArray []string

	body = strings.Replace(body, "\n", "", -1)

	// parsing title and year
	re := regexp.MustCompile(IMDB_TITLE)
	titleAndYearRaw := re.FindAllString(body, -1)

	if titleAndYearRaw == nil {
		fmt.Println("No matches.")
	} else {
		r, _ := regexp.Compile("[0-9]{4}")
		yearIndex := r.FindStringIndex(titleAndYearRaw[0])
		if yearIndex == nil {
			fmt.Println("No matches.")
		} else {
			r, _ = regexp.Compile("=\"")
			titleIndex := r.FindStringIndex(titleAndYearRaw[0])
			if titleIndex == nil && len(titleIndex) < 2 {
				fmt.Println("No matches.")
			} else {
				year, _ = strconv.ParseInt(titleAndYearRaw[0][yearIndex[0]:yearIndex[1]], 10, 32)
				title = titleAndYearRaw[0][titleIndex[0]+2 : yearIndex[0]-2]
			}
		}
	}
	// parsing rating
	re = regexp.MustCompile(IMDB_RATING)
	ratingRawString := re.FindAllString(body, -1)

	if ratingRawString == nil {
		fmt.Println("No matches.")
	} else {
		r, _ := regexp.Compile("[+-]?([0-9]*[.])?[0-9]+")
		ratingStrings := r.FindAllString(ratingRawString[0], -1)
		if ratingStrings == nil {
			fmt.Println("No matches.")
		} else {
			rating, _ = strconv.ParseFloat(ratingStrings[0], 8)
		}
	}
	// parsing summary
	re = regexp.MustCompile(IMDB_SUMMARY)
	summaryRaw := re.FindAllString(body, -1)

	if summaryRaw == nil {
		fmt.Println("No matches.")
	} else {
		r, _ := regexp.Compile(">(.*?)<")
		summaryIndex := r.FindStringIndex(summaryRaw[0])
		if summaryIndex == nil || len(summaryIndex) < 2 {
			fmt.Println("No matches.")
		} else {
			summary = summaryRaw[0][summaryIndex[0]+1 : summaryIndex[1]-1]
			summary = strings.TrimRight(summary, " ")
			summary = strings.TrimLeft(summary, " ")
		}
	}
	// parsing duration
	re = regexp.MustCompile(IMDB_DURATION)
	durationRaw := re.FindAllString(body, -1)

	if durationRaw == nil {
		fmt.Println("No matches.")
	} else {
		r, _ := regexp.Compile("<time(.*?)</time>")
		if r == nil {
			fmt.Println("No matches")
		} else {
			durationString := r.FindAllString(durationRaw[0], -1)
			if durationString == nil {
				fmt.Println("No matches.")
			} else {
				r, _ = regexp.Compile(">(.*?)<")
				durationIndex := r.FindStringIndex(durationString[0])

				if len(durationIndex) < 2 {
					fmt.Println("No matches.")
				} else {
					duration = durationString[0][durationIndex[0]+1 : durationIndex[1]-1]
				}
			}
		}
	}
	// matching genre
	re = regexp.MustCompile(IMDB_GENRES)
	genresRaw := re.FindAllString(body, -1)

	if genresRaw == nil {
		fmt.Println("No matches.")
	} else {
		r, _ := regexp.Compile("<a(.*?)</a>")
		genreRaw := r.FindAllString(genresRaw[0], -1)
		if genreRaw == nil {
			fmt.Println("No matches.")
		} else {
			for _, genre := range genreRaw {
				r, _ = regexp.Compile(">(.*?)<")
				genreIndex := r.FindStringIndex(genre)
				genre := genre[genreIndex[0]+1 : genreIndex[1]-1]
				genreArray = append(genreArray, genre)
			}
		}
	}
	return &Movie{
		Title:    title,
		Year:     year,
		Rating:   rating,
		Summary:  summary,
		Duration: duration,
		Genre:    strings.Join(genreArray, ","),
	}

}
