package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	htmlFiles, err := getAllHTMLFiles("./cmd/web")

	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.LoadHTMLFiles(htmlFiles...)

	router.GET("/room/:roomID", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat-room.html", nil)
	})

	router.GET("/ws/:roomID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		fmt.Println(roomID)
	})

	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatal(err)
	}
}
