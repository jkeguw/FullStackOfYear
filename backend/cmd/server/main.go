package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jkeguw/FullStackOfYear/internal/api/routes"
	"github.com/jkeguw/FullStackOfYear/pkg/common/config"
	"log"
)

func main() {
	config.Init()
	r := gin.Default()
	routes.SetupRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
