package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wbergg/beerel-roll/db"
)

var ProductNumber struct {
	Pid      int64  `json:"pid"`
	Consumer string `json:"consumer"`
}

func Start(d *db.DBobject) {

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
			err := d.ConsumeBeer(data.Pid, data.Consumer)
			if err != nil {
				log.Fatalf("Could not remove from DB: %v", err)
			}
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			log.Fatalf("DBERR: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// Endpoint to reject a roll
	r.POST("/reject", func(c *gin.Context) {
		// Declare a variable of the struct type
		data := ProductNumber
		// Bind JSON request body to the variable
		if err := c.ShouldBindJSON(&data); err == nil {
			// Process the data
			err := d.RejectBeer(data.Pid, data.Consumer)
			if err != nil {
				log.Fatalf("Could not remove from DB: %v", err)
			}
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			log.Fatalf("DBERR: %v", err)
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
