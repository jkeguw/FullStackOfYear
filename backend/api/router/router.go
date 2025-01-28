package router

import (
	"FullStackOfYear/backend/api/v1"
	"FullStackOfYear/backend/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	apiV1 := r.Group("/api/v1")
	v1.RegisterRoutes(apiV1)

	return r
}
