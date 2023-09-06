package routers

import (
	"examn_go/controllers"
	"examn_go/routers/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })
	route.POST("/auth/login", controllers.LoginHandler)
	route.GET("/:exam_id/questions", controllers.GetAllQuestions)
	route.POST("/tests/:exam_id/start",middleware.AuthMiddleware(), controllers.StartTesHandler)
	route.GET("/tests/attempt/ongoing",middleware.AuthMiddleware(), controllers.GetOngoingAttempt)
	route.POST("/tests/attempt/:attempt_id/userevent", controllers.PostUserEvents)
	route.GET("/tests/attempt/:attempt_id/reload", controllers.ReloadAttempt)
	route.POST("/tests/attempt/:attempt_id/end", controllers.EndTest)
	route.POST("/auth/refresh", controllers.RefreshTokensHandler)


	//Add All route
	//TestRoutes(route)
}
