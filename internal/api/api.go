package api

import (
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru"
	"github.com/tazer/ipwhois/internal/database"
)

type Api struct {
	database *database.Database
	Web      *gin.Engine
	cache    *lru.Cache
}

func NewApi(database *database.Database, cache *lru.Cache) *Api {
	api := &Api{
		Web:      gin.Default(),
		database: database,
		cache:    cache,
	}
	api.setup()

	return api
}

func (a *Api) setup() {
	a.Web.GET("/", a.startHandler)
	a.Web.GET("/ip/:ip", a.ipwhoisHandler)
}

func (a *Api) startHandler(c *gin.Context) {
	c.String(200, fmt.Sprintf(`This product includes GeoLite2 data created by MaxMind, available from
	<a href="https://www.maxmind.com">https://www.maxmind.com</a>. \n \n
	Try it out here: GET /ip/%s`, c.ClientIP()), nil)
}

func (a *Api) ipwhoisHandler(c *gin.Context) {
	ip := c.Param("ip")

	if cachedCountry, ok := a.cache.Get(ip); ok {
		etag := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", cachedCountry))))
		c.Header("ETag", etag)
		c.Header("Cache-Control", "public, max-age=86400")
		if match := c.GetHeader("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				c.Status(http.StatusNotModified)
				return
			}
		}
		c.JSON(200, cachedCountry)
		return
	}
	realIP := net.ParseIP(ip)

	country, err := a.database.DB.Country(realIP)

	if err != nil {
		log.Printf("Error gettign country err: %v", err)
		c.AbortWithError(500, err)
		return
	}
	a.cache.Add(ip, country)
	etag := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v", country))))
	c.Header("ETag", etag)
	c.Header("Cache-Control", "public, max-age=86400")
	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			c.Status(http.StatusNotModified)
			return
		}
	}
	c.JSON(200, country)
}
