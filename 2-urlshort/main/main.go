package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"com.elpigo/urlshort"
	"github.com/boltdb/bolt"
)

var filename string
var boltDB *bool
var boltDbPath string = "urlshort.db"
var handler string

func main() {

	//flag options
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information: \n")
		flag.PrintDefaults()
	}

	flag.StringVar(&filename, "f", "", "use -f to pass dir path and filename of YAML file")
	boltDB = flag.Bool("bolt", false, "use flag -bolt to specify BoltDB as source of key/value pairs")
	path := flag.String("path", "", "add short url path, use following format: /path")
	url := flag.String("url", "", "full url to map, use following format: http://foo.com")
	view := flag.Bool("view", false, "use flag in conjunctin with flag -b bolt to view current key/value entries in Bolt")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/google":         "https://www.google.com",
	}
	mux := defaultMux()
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback. Hardcode some yaml paths if no yaml-file is present
	var yaml string
	if filename == "" {
		yaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
- path: /fb
  url: https://www.facebook.com
  `
	} else {
		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Could not read file: %s due to %v\n", filename, err)
		}
		yaml = string(yamlFile)
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	handler = "yamlHandler"

	//if using boltDB
	var boltHandler http.HandlerFunc
	if *boltDB {
		db := useBoltDb()
		defer db.Close()

		if *path != "" && *url != "" {
			urlshort.AddEntryBoltDb(db, *path, *url)
		}

		//call BoltDB handler
		boltHandler, err = urlshort.BoltDBHandler(db, yamlHandler)
		if err != nil {
			log.Println(err)
		}

		if *view {
			urlshort.ViewEntriesBoltDb(db)
		}

		handler = "boltHandler"

	}

	//start the server. Used a switch statement as cleaner to keep things together here.
	switch handler {
	case "yamlHandler":
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", yamlHandler)
	case "boltHandler":
		fmt.Println("Starting the server on :8080")
		http.ListenAndServe(":8080", boltHandler)
	}

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "No route configured for this endpoint")
}

func useBoltDb() *bolt.DB {
	db, err := urlshort.BoltDBSetup()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
