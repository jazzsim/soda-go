package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type HttpServer struct {
	Url string `json:"url"`
}

func main() {
	router := gin.Default()
	router.POST("/scrape", scrape)

	router.Run("localhost:8080")
}

func scrape(c * gin.Context) {
	server := HttpServer{}

	if err := c.ShouldBindJSON(&server); err != nil {
		// Handle error (e.g., invalid JSON format)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Instantiate default collector
	collector := colly.NewCollector()

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	contents := Contents{}
	contents.retrieveContents(collector)

	collector.OnScraped(func(r *colly.Response) {
 	   fmt.Println("Finished", r.Request.URL)
		c.IndentedJSON(http.StatusOK, contents)
	})

	collector.Visit(server.Url)
}