package api

import (
	"github.com/Kurlabs/alerty/cmd/brain/api/handlers"
	"github.com/labstack/echo"
)

func JwtGroup(g *echo.Group) {
	g.GET("/check", handlers.MainJwt)
	g.POST("/monitors/batch", handlers.MonitorBatch)
	g.POST("/monitors/robot", handlers.MonitorRobot)
}
