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
	// PassengerService is used to handle the passenger use cases
	PassengerService interface {
		Create(*models.Passenger, *auth.UserClaims) error
		Get(id string, user *auth.UserClaims) (*models.Passenger, error)
		Update(updates *models.PassengerChangeSet, passengerID string, user *auth.UserClaims) (*models.Passenger, error)
		GetAll(lastID string, limit, userAuthLevel int) ([]*models.Passenger, error)
		GetAllByRideID(rideID string, user *auth.UserClaims) ([]*models.Passenger, error)
		Delete(id string, user *auth.UserClaims) error
	}

	// PassengerController http adapter
	PassengerController struct {
		logger  interfaces.Logger
		service PassengerService
	}
)

// NewPassengerController creates a new passenger controller
func NewPassengerController(r PassengerService, l interfaces.Logger) *PassengerController {
	return &PassengerController{
		logger:  l,
		service: r,
	}
}

// MountRoutes mounts the auth routes
func (p *PassengerController) MountRoutes(c *echo.Group) {
	c.GET("/passengers", p.list, auth.NewAuthMiddleware(services.AdminLevel, p.logger))
	c.GET("/passengers/:id", p.show)
	c.Use(auth.NewAuthMiddleware(services.UserLevel, p.logger))
	c.POST("/passengers", p.create)
	c.DELETE("/passengers/:id", p.delete)
	c.PUT("/passengers/:id", p.update)
}

func (p *PassengerController) create(c echo.Context) error {
	data := models.PassengerChangeSet{}
	if err := c.Bind(&data); err != nil {
		msg := err.Error()
		if strings.HasPrefix(err.Error(), "code=400, message=Syntax error") {
			msg = "The JSON was invalid"
		}

		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	user := userClaimsFromContext(c)
	data.PassengerID = &user.ID
	status := models.PassengerInterested
	data.Status = &status

	passenger, err := models.NewPassenger(&data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := p.service.Create(passenger, user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": passenger,
	})
}

func (p *PassengerController) list(c echo.Context) error {
	user := userClaimsFromContext(c)
	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)

	if limitStr == "" {
		limit = 15
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "limit must be an integer greater than 0")
	}

	lastID := c.QueryParam("last")

	passengers, err := p.service.GetAll(lastID, limit, user.AuthLevel)

	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrNotAllowed {
			status = http.StatusForbidden
		}

		return echo.NewHTTPError(status, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": passengers,
	})
}

func (p *PassengerController) show(c echo.Context) error {
	id := c.Param("id")
	user := userClaimsFromContext(c)

	passenger, err := p.service.Get(id, user)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": passenger,
	})
}

func (p *PassengerController) update(c echo.Context) error {
	id := c.Param("id")

	data := models.PassengerChangeSet{}
	if err := c.Bind(&data); err != nil {
		msg := err.Error()
		if strings.HasPrefix(err.Error(), "code=400, message=Syntax error") {
			msg = "The JSON was invalid"
		}

		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	// only allow status updates
	data.RideID = nil
	data.PassengerID = nil

	user := userClaimsFromContext(c)

	passenger, err := p.service.Update(&data, id, user)

	if err != nil {
		status := http.StatusBadRequest
		if err == services.ErrForbidden {
			status = http.StatusForbidden
		}

		return echo.NewHTTPError(status, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": passenger,
	})
}

func (p *PassengerController) delete(c echo.Context) error {
	id := c.Param("id")

	user := userClaimsFromContext(c)

	err := p.service.Delete(id, user)

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
