package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type val struct {
	id int
	s  string
}

var (
	static  = flag.String("staticdir", "static", "the directory storing the static web files")
	csvPath = flag.String("csv", "", "the input csv filepath")
	lonField = flag.String("lon", "", "the field containing longitude values")
	latField = flag.String("lat", "", "the field containing latitude values")
)

func main() {
	header, c, data := setup()
	fmt.Println("done setting up")

	http.Handle("/", http.FileServer(http.Dir(*static)))
	http.Handle("/getheader/", headerHandler(header))
	http.Handle("/search/", searchHandler(c, data))

	fmt.Println("serving at localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func headerHandler(header []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(header)
		if err != nil {
			fmt.Fprintf(w, fmt.Sprintf("%v", err))
		}
		w.Write(j)
	})
}

func searchHandler(c [][]string, data []val) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("%v", err))
			}
			s := r.FormValue("search")
			log.Println(s)

			s = standardize(s)

			rows, err := findMatchingRows(c, data, s)
			if err != nil {
				j, _ := json.Marshal(err.Error())
				w.Write(j)
				return
			}

			j, err := json.Marshal(rows)
			if err != nil {
				fmt.Fprintf(w, "server error")
			}
			w.Write(j)
		}
	})
}

func findMatchingRows(c [][]string, data []val, s string) ([][]string, error) {
	i := sort.Search(len(data), func(i int) bool {
		return data[i].s >= s
	})

	if data[i].s == s {
		keys := map[int]bool{}

		foundId := data[i].id
		rows := [][]string{c[foundId]}
		keys[foundId] = true

		for j := 1; data[i+j].s == s; j++ {
			foundId = data[i+j].id
			if _, found := keys[foundId]; !found {
				rows = append(rows, c[foundId])
				keys[foundId] = true
			}
		}

		return rows, nil
	} else {
		return [][]string{}, fmt.Errorf("not found")
	}
}

func setup() ([]string, [][]string, []val) {
	flag.Parse()

	f, err := os.Open(*csvPath)
	if err != nil {
		panic(err)
	}

	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		panic(err)
	}

	hasLocation :=  *lonField != "" && *latField != ""
	var (
		lonIndex = -1
		latIndex = -1
	)

	if hasLocation {
		// recreate header
		newHeader := []string{}
		for i, field := range header {
			// if it's anything other than lon/lat, add it to the new header
			if field != *lonField && field != *latField {
				newHeader = append(newHeader, field)
			} else {
				// else find indexes
				if field == *lonField {
					lonIndex = i
				}
				if field == *latField {
					latIndex = i
				}
			}
		}
		header = newHeader
	}

	lookup := []val{}
	c := [][]string{}

	i := 0
	for  {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		var lon, lat string
		cRow := []string{}
		for j, v := range row {
			if j != lonIndex && j != latIndex {
				lookup = append(lookup, val{id: i, s: standardize(v)})
				cRow = append(cRow, v)
			}
			
			// Create point from data
			if j == lonIndex {
				lon = v
			}
			if j == latIndex {
				lat = v
			}
		}
		if hasLocation {
			link := constructLink(lon, lat)
			cRow = append(cRow, link)
		}

		i++
		c = append(c, cRow)
	}

	sort.Slice(lookup, func(i, j int) bool {
		return lookup[i].s < lookup[j].s
	})

	return header, c, lookup
}

func constructLink(lon, lat string) string {
	return fmt.Sprintf(
		"https://www.google.com/maps/search/?api=1&query=%s,%s&authuser=1",
		lat, lon,
	)
}

func standardize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
