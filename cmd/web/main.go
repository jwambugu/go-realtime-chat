package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	htmlFiles, err := getAllHTMLFiles("./cmd/web")

	if err != nil {
		log.Fatal(err)
	}

	go h.run()

	router := gin.New()
	router.LoadHTMLFiles(htmlFiles...)

	router.GET("/room/:roomID", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/ws/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")

		serveWebsockets(c.Writer, c.Request, roomID)
	})

	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}
}
