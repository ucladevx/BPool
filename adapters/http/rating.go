package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/utils/auth"
)

type (
	// RatingService is used to provide and check ratings
	RatingService interface {
		Create(models.Rating, *auth.UserClaims) (*models.Rating, error)
		GetByID(string, *auth.UserClaims) (*models.Rating, error)
		GetRatingByUserID(string) (float32, error)
		Delete(string, *auth.UserClaims) error
	}

	// RatingController http adapter
	RatingController struct {
		service          RatingService
		passengerService PassengerService
		logger           interfaces.Logger
	}
)

// NewRatingController creates a new rating controller
func NewRatingController(ratingService RatingService, p PassengerService, l interfaces.Logger) *RatingController {
	return &RatingController{
		service:          ratingService,
		passengerService: p,
		logger:           l,
	}
}

// MountRoutes mounts the rating routes
func (ratingController *RatingController) MountRoutes(c *echo.Group) {
	c.DELETE(
		"/ratings/:id",
		ratingController.delete,
		auth.NewAuthMiddleware(services.AdminLevel, ratingController.logger),
	)

	c.Use(auth.NewAuthMiddleware(services.UserLevel, ratingController.logger))

	c.GET("/ratings/:id", ratingController.show)
	c.GET("/ratings/user/:id", ratingController.getUserRating)
	c.POST("/ratings", ratingController.create)
}

func (ratingController *RatingController) show(c echo.Context) error {
	id := c.Param("id")
	user := userClaimsFromContext(c)

	rating, err := ratingController.service.GetByID(id, user)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": rating,
	})
}

func (ratingController *RatingController) getUserRating(c echo.Context) error {
	userID := c.Param("id")

	rating, err := ratingController.service.GetRatingByUserID(userID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": rating,
	})
}

func (ratingController *RatingController) create(c echo.Context) error {
	data := models.Rating{}

	if err := c.Bind(&data); err != nil {
		message := err.Error()
		if strings.HasPrefix(message, "code=400, message=Syntax error") {
			message = "Invalid JSON"
		}

		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	userClaims := userClaimsFromContext(c)
	rating, err := ratingController.service.Create(data, userClaims)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": rating,
	})
}

func (ratingController *RatingController) delete(c echo.Context) error {
	id := c.Param("id")
	userClaims := userClaimsFromContext(c)

	err := ratingController.service.Delete(id, userClaims)

	if err != nil {
		status := 400
		// fill in other errors

		return echo.NewHTTPError(status, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
