package router

import (
	"github.com/kamalshkeir/muzzsol/handlers"
	"github.com/kamalshkeir/muzzsol/middlewares"
	"github.com/labstack/echo/v4"
)

func InitUrls(e *echo.Echo) {
	// it make sense to make it as POST request, but here for convenience of demonstration in the browser and because we don't need to receive body data, i keep it as GET request
	e.GET("/user/create",handlers.UserRandom)
	e.GET("/profiles",middlewares.AuthMiddleware(handlers.GetUserProfiles) )
	e.POST("/swipe",middlewares.AuthMiddleware(handlers.UserSwipe))
	e.POST("/login",handlers.Login)
}