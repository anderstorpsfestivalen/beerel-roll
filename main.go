package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wbergg/beerel-roll/dataimport"
	"github.com/wbergg/beerel-roll/db"
)

var ProductNumber struct {
	Pid int64 `json:"pid"`
}

func main() {
	// Setup db
	d := db.Open()

	// Check if DB is set up, if not, set it up (first time only)
	if d.Setup == 0 {
		fmt.Println("Looks like it's the first time - Populating DB...")
		err := dataimport.DbSetup(&d)
		if err == nil {
			fmt.Println("DB population sucess! Please rerun the program!")
		}
		os.Exit(0)
	}

	// Setup Gin
	r := gin.Default()

	// Serve the HTML page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Endpoint to get a random beer
	r.GET("/random-beer", func(c *gin.Context) {
		beer, err := d.GetRandBeer()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, beer)
		}
	})

	// Endpoint to accept a roll
	r.POST("/accept", func(c *gin.Context) {
		// Declare a variable of the struct type
		data := ProductNumber

		// Bind JSON request body to the variable
		if err := c.ShouldBindJSON(&data); err == nil {
			// Process the data
			err := d.ConsumeBeer(data.Pid)
			if err != nil {
				log.Fatalf("Could not remove from DB: %v", err)
			}
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// Endpoint to get a 5 last beers consumed
	r.GET("/recent", func(c *gin.Context) {
		beer, err := d.GetNLastConsumed(5)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, beer)
		}
	})

	// Load the HTML template
	r.LoadHTMLFiles("templates/index.html")

	// Start webserver
	r.Run(":8080")
}
