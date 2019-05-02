package ipwhois

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/mholt/archiver"
)

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

func GetFile() error {
	err := getMaxMindArchive()

	if err != nil {
		return err
	}

	err = unpackArchive()

	if err != nil {
		return err
	}

	return nil
}

func unpackArchive() error {
	_, err := os.Stat("GeoIP2-Country.mmdb")

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
						return err
					}
					nf.Write(buf)

				}
				return nil
			})

			if err != nil {
				return err
			}
		} else {
			return err
		}

	}
	return nil
}

func getMaxMindArchive() error {
	_, err := os.Stat("./GeoLite2-Country.tar.gz")

	// check if there is an error
	if err != nil {

		// check if error is file does not exist
		if os.IsNotExist(err) {
			downloadFile("./GeoLite2-Country.tar.gz", "https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.tar.gz")
			_, err = os.Stat("./GeoLite2-Country.tar.gz")
			if err != nil {
				return err
			}
		} else {
			return err
		}

	}

	// TODO: Implement cron that checks every 24h for new file
	// Implement m5 hash of the files so it can be compared against
	// https://geolite.maxmind.com/download/geoip/database/GeoLite2-Country.md5
	// If there is a new file download and reopen database with the new file

	return nil
}
