package http

import (
	"net/http"

	"github.com/labstack/echo"
)

type Pages struct {
	logger Logger
}

func NewPagesController(l Logger) *Pages {
	return &Pages{
		logger: l,
	}
}

func (p *Pages) MountRoutes(c *echo.Group) {
	c.GET("/", p.index)
	c.GET("/health", p.health)
}

func (p *Pages) index(c echo.Context) error {
	return c.HTML(http.StatusOK, "<html><title>Golang Google</title> <body> <a href='/api/auth/login'><button>Login with Google!</button> </a> </body></html>")
}

func (p *Pages) health(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"health": "OK",
	})
}
