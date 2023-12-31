package main

import (
	"fmt"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type HttpServer struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.POST("/scrape", scrape)

	router.Run("localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func scrape(c *gin.Context) {
	server := HttpServer{}

	if err := c.ShouldBindJSON(&server); err != nil {
		// Handle error (e.g., invalid JSON format)
		c.JSON(http.StatusBadRequest, gin.H{"400 error": err.Error()})
		return
	}

	_, err := server.verifyUrl()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"500 error": err.Error()})
		return
	}

	// Instantiate default collector
	collector := colly.NewCollector()

	credentialsNeeded := server.Username != "" && server.Password != ""
	if credentialsNeeded {
		fmt.Println("in credentials")
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(server.Username+":"+server.Password))
		collector.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Authorization", auth)
		})
	}

	contents := Contents{Folders: []string{}, Files: []Files{}}
	contents.retrieveContents(collector)

	collector.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		c.IndentedJSON(http.StatusOK, contents)
	})

	collector.Visit(server.Url)
}

func (s *HttpServer) verifyUrl() (*url.URL, error) {
	// Parse the URL
	parseUrl, err := url.Parse(s.Url)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil, err
	}
	return parseUrl, nil
}
