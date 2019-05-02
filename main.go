package main

import (
	"log"

	lru "github.com/hashicorp/golang-lru"
	"github.com/oschwald/geoip2-golang"
	"github.com/tazer/ipwhois/internal/api"
	"github.com/tazer/ipwhois/internal/database"
	"github.com/tazer/ipwhois/internal/ipwhois"
)

func main() {

	err := ipwhois.GetFile()

	if err != nil {
		log.Fatal(err)
	}

	lruCache, err := lru.New(1024)

	if err != nil {
		log.Fatalf("Coudlnt create LRU cache err: %v", err)
	}

	db, err := geoip2.Open("GeoIP2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	database := database.NewDatabase(db)

	api := api.NewApi(database, lruCache)
	api.Web.Run() // listen and serve on 0.0.0.0:8080
}

// func hash_file_md5(filePath string) (string, error) {
// 	//Initialize variable returnMD5String now in case an error has to be returned
// 	var returnMD5String string

// 	//Open the passed argument and check for any error
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return returnMD5String, err
// 	}

// 	//Tell the program to call the following function when the current function returns
// 	defer file.Close()

// 	//Open a new hash interface to write to
// 	hash := md5.New()

// 	//Copy the file in the hash interface and check for any error
// 	if _, err := io.Copy(hash, file); err != nil {
// 		return returnMD5String, err
// 	}

// 	//Get the 16 bytes hash
// 	hashInBytes := hash.Sum(nil)[:16]

// 	//Convert the bytes to a string
// 	returnMD5String = hex.EncodeToString(hashInBytes)

// 	return returnMD5String, nil

// }
