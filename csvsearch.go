package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
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
		rows := [][]string{c[data[i].id]}

		for j := 1; data[i+j].s == s; j++ {
			rows = append(rows, c[data[i+j].id])
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

	c, err := r.ReadAll()
	if err != nil {
		panic(err)
	}

	data := []val{}

	for i, row := range c {
		if err != nil {
			break
		}
		for _, v := range row {
			if v != "" {
				data = append(data, val{id: i, s: standardize(v)})
			}
		}
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].s < data[j].s
	})

	return header, c, data
}

func standardize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
