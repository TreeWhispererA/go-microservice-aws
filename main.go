package main

import (
	"fmt"
	"log"
	"os"

	"blockparty.co/test/routes"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	fmt.Println("Port number is: " + port)
	color.Cyan("🌏 Server running on localhost:" + port)

	router := gin.Default()
	router.Use(gin.Logger())
	routes.Routes(router)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
