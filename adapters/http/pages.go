package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/interfaces"
)

// PagesController is the controller for pages
type PagesController struct {
	logger interfaces.Logger
}

// NewPagesController creates a new pages controller
func NewPagesController(l interfaces.Logger) *PagesController {
	return &PagesController{
		logger: l,
	}
}

// MountRoutes adds the pages routes to the apps
func (p *PagesController) MountRoutes(c *echo.Group) {
	c.GET("/", p.index)
	c.GET("/health", p.health)
}

func (p *PagesController) index(c echo.Context) error {
	return c.String(http.StatusOK, "Hello there!")
}

func (p *PagesController) health(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"health": "OK",
	})
}
