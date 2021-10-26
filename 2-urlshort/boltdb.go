package urlshort

import (
	"fmt"

	"github.com/boltdb/bolt"
)

var dbBucket string = "UrlShort"
var boltDbName string = "urlshort.db"

type BoltData struct {
	Path string
	Url  string
}

func BoltDBSetup() (*bolt.DB, error) {
	db, err := bolt.Open(boltDbName, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not open db %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(dbBucket))
		if err != nil {
			return fmt.Errorf("could not create BoltDB bucket: %v", err)
		}
		return nil
	})
	fmt.Println("Using BoltDB")
	return db, nil
}

func AddEntryBoltDb(db *bolt.DB, path string, url string) error {

	entry := BoltData{Path: path, Url: url}

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(dbBucket)).Put([]byte(entry.Path), []byte(entry.Url))
		if err != nil {
			return fmt.Errorf("Could not insert path and url into BoltDB: %v", err)
		}
		return nil
	})
	fmt.Printf("Added path: %s and url: %s\n", entry.Path, entry.Url)
	return err
}

func getEntriesBoltDb(db *bolt.DB) (map[string]string, error) {

	entries := make(map[string]string)
	err := db.View(func(tx *bolt.Tx) error {
		values := tx.Bucket([]byte("UrlShort"))
		values.ForEach(func(k, v []byte) error {
			entries[string(k)] = string(v)
			return nil
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}

//ViewEntriesBoltDb will show entries in BoltDb in YAML format
func ViewEntriesBoltDb(db *bolt.DB) {
	entries, err := getEntriesBoltDb(db)
	if err != nil {
		fmt.Println(err)
	}
	toYaml(entries)
}
