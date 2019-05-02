package database

import "github.com/oschwald/geoip2-golang"

type Database struct {
	DB *geoip2.Reader
}

func NewDatabase(db *geoip2.Reader) *Database {
	return &Database{
		DB: db,
	}
}
