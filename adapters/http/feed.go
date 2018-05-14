package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/utils/auth"
)

type (
	// FeedService handles the use cases for the feed
	FeedService interface {
		GetUserRides(userID string) ([]*models.Ride, error)
	}

	// FeedController is the controller for the feed
	FeedController struct {
		logger      interfaces.Logger
		feedService FeedService
	}
)

// NewFeedController creates a new feed controller
func NewFeedController(f FeedService, l interfaces.Logger) *FeedController {
	return &FeedController{
		logger:      l,
		feedService: f,
	}
}

// MountRoutes adds the pages routes to the apps
func (f *FeedController) MountRoutes(c *echo.Group) {
	c.Use(auth.NewAuthMiddleware(services.UserLevel, f.logger))
	c.GET("/feed", f.getFeed)
}

func (f *FeedController) getFeed(c echo.Context) error {
	user := userClaimsFromContext(c)

	rides, err := f.feedService.GetUserRides(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": rides,
	})
}
