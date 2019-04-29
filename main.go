package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver"
	"github.com/oschwald/geoip2-golang"
)

func main() {

	//https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.tar.gz

	_, err := os.Stat("GeoLite2-Country.tar.gz")

	// check if there is an error
	if err != nil {

		// check if error is file does not exist
		if os.IsNotExist(err) {
			downloadFile("GeoLite2-Country.tar.gz", "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.tar.gz")
			_, err = os.Stat("GeoLite2-Country.tar.gz")
			if err != nil {
				log.Fatalf("Coudlnt download file info err: %v", err)
			}
		} else {
			log.Fatalf("Coudlnt read file info err: %v", err)
		}

	}

	_, err = os.Stat("GeoIP2-Country.mmdb")

	// check if there is an error
	if err != nil {

		// check if error is file does not exist
		if os.IsNotExist(err) {
			err = archiver.Walk("GeoLite2-Country.tar.gz", func(f archiver.File) error {
				if strings.Contains(f.Name(), ".mmdb") {
					nf, err := os.Create("GeoIP2-Country.mmdb")
					defer nf.Close()
					buf, err := ioutil.ReadAll(f)
					if err != nil {
						log.Fatal(err)
					}
					nf.Write(buf)

				}
				return nil
			})

			if err != nil {
				log.Fatalf("Error unarchive err: %v", err)
			}
		} else {
			log.Fatalf("Coudlnt read file info err: %v", err)
		}

	}

	//TODO: Implement hash

	db, err := geoip2.Open("GeoIP2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, fmt.Sprintf(`This product includes GeoLite2 data created by MaxMind, available from
		<a href="https://www.maxmind.com">https://www.maxmind.com</a>. \n \n
		Try it out here: GET /ip/%s`, c.ClientIP()), nil)
	})
	r.GET("/ip/:ip", func(c *gin.Context) {
		ip := c.Param("ip")

		realIp := net.ParseIP(ip)

		country, err := db.Country(realIp)

		if err != nil {
			log.Printf("Error gettign country err: %v", err)
			c.AbortWithError(500, err)
			return
		}

		c.JSON(200, country)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func hash_file_md5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}