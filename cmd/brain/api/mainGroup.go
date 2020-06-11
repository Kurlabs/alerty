package api

import (
	"github.com/Kurlabs/alerty/cmd/brain/api/handlers"
	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	// e.GET("/login", handlers.Login)

	e.GET("/", handlers.Index)
}
