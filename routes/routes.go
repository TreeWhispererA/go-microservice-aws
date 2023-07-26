package routes

import (
	"blockparty.co/test/controllers"
	"github.com/gin-gonic/gin"
)

// Routes -> define endpoints
func Routes(router *gin.Engine) {

	router.GET("/test", controllers.Test)

	router.GET("/tokens", controllers.GetTokens)
	router.GET("/tokens/:cid", controllers.GetTokenByID)
	router.GET("/scrap", controllers.ScrapFunc)
}
