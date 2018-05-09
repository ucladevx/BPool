package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/utils/auth"
)

type (
	// RideService is used to handle the ride use cases
	RideService interface {
		Create(*models.Ride, *auth.UserClaims) error
		Get(id string) (*models.Ride, error)
		Update(updates *models.RideChangeSet, rideID string, user *auth.UserClaims) (*models.Ride, error)
		GetAll(lastID string, limit, userAuthLevel int) ([]*models.Ride, error)
		Delete(id string, user *auth.UserClaims) error
	}

	// RideController http adapter
	RideController struct {
		logger           interfaces.Logger
		service          RideService
		passengerService PassengerService
	}
)

// NewRideController creates a new auth controller
func NewRideController(r RideService, p PassengerService, l interfaces.Logger) *RideController {
	return &RideController{
		logger:           l,
		service:          r,
		passengerService: p,
	}
}

// MountRoutes mounts the auth routes
func (r *RideController) MountRoutes(c *echo.Group) {
	c.GET("/rides", r.list, auth.NewAuthMiddleware(services.AdminLevel, r.logger))
	c.GET("/rides/:id", r.show)
	c.Use(auth.NewAuthMiddleware(services.UserLevel, r.logger))
	c.POST("/rides", r.create)
	c.DELETE("/rides/:id", r.delete)
	c.PUT("/rides/:id", r.update)
}

func (r *RideController) create(c echo.Context) error {
	data := models.RideChangeSet{}
	if err := c.Bind(&data); err != nil {
		msg := err.Error()
		if strings.HasPrefix(err.Error(), "code=400, message=Syntax error") {
			msg = "The JSON was invalid"
		}

		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	user := userClaimsFromContext(c)
	data.DriverID = &user.ID

	ride, err := models.NewRide(&data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := r.service.Create(ride, user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": ride,
	})
}

func (r *RideController) list(c echo.Context) error {
	user := userClaimsFromContext(c)
	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)

	if limitStr == "" {
		limit = 15
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "limit must be an integer greater than 0")
	}

	lastID := c.QueryParam("last")

	rides, err := r.service.GetAll(lastID, limit, user.AuthLevel)

	if err != nil {
		if err == services.ErrNotAllowed {
			return ErrNotAllowed
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": rides,
	})
}

func (r *RideController) show(c echo.Context) error {
	id := c.Param("id")

	ride, err := r.service.Get(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": ride,
	})
}

func (r *RideController) update(c echo.Context) error {
	id := c.Param("id")

	data := models.RideChangeSet{}
	if err := c.Bind(&data); err != nil {
		msg := err.Error()
		if strings.HasPrefix(err.Error(), "code=400, message=Syntax error") {
			msg = "The JSON was invalid"
		}

		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	user := userClaimsFromContext(c)

	ride, err := r.service.Update(&data, id, user)

	if err != nil {
		status := http.StatusBadRequest
		if err == services.ErrForbidden {
			status = http.StatusForbidden
		}

		return echo.NewHTTPError(status, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": ride,
	})
}

func (r *RideController) delete(c echo.Context) error {
	id := c.Param("id")

	user := userClaimsFromContext(c)

	err := r.service.Delete(id, user)

	if err != nil {
		status := 400
		if err == services.ErrNotAllowed {
			status = 401
		} else if err == services.ErrForbidden {
			status = 403
		}

		return echo.NewHTTPError(status, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
