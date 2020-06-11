package router

import (
	"github.com/Kurlabs/alerty/cmd/brain/api"
	"github.com/Kurlabs/alerty/cmd/brain/api/middlewares"
	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// create groups
	jwtGroup := e.Group("/api")

	// set all middlewares
	middlewares.SetJwtMiddlewares(jwtGroup)

	// set routes
	api.MainGroup(e)
	api.JwtGroup(jwtGroup)

	return e
}
