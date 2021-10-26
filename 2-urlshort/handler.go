package urlshort

import (
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if val, ok := pathsToUrls[r.URL.Path]; ok {
			log.Printf("redirecting path: %s to: %s", r.URL.Path, val)
			http.Redirect(w, r, val, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}

}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := ParseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := BuildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

//BoltDBHandler will take the key which is defined as the api path, look it up in
//BoltDB and return the associated value. If it doesn't exist it will fallback onto the
//http.Hanlder which just returns a hello-world default page. If it does exist, it will marshal the value
//into YAML, pass it to BuildMap and then pass it to MapHandler
func BoltDBHandler(db *bolt.DB, fallback http.Handler) (http.HandlerFunc, error) {
	entries, err := getEntriesBoltDb(db)
	if err != nil {
		return nil, err
	}
	return MapHandler(entries, fallback), nil
}
