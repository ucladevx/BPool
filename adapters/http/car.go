package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/stores/postgres"
	"github.com/ucladevx/BPool/utils/auth"
)

type (
	// CarService is used to handle all car CRUD operations
	CarService interface {
		GetAll(lastID string, limit, authLevel int) ([]*models.Car, error)
		GetCar(id string) (*models.Car, error)
		AddCar(body models.Car, userID string) (*models.Car, error)
		DeleteCar(id, userID string) error
	}

	// CarController http adapter
	CarController struct {
		logger  interfaces.Logger
		service CarService
	}
)

// NewCarController creates a new car controller
func NewCarController(c CarService, l interfaces.Logger) *CarController {
	return &CarController{
		logger:  l,
		service: c,
	}
}

// MountRoutes mounts the car routes
func (cc *CarController) MountRoutes(c *echo.Group) {
	c.GET("/cars", cc.list, auth.NewAuthMiddleware(services.AdminLevel, cc.logger))

	c.Use(auth.NewAuthMiddleware(services.UserLevel, cc.logger))

	c.GET("/cars/:id", cc.show)
	c.POST("/cars", cc.create)
	c.DELETE("/cars/:id", cc.remove)
}

func (cc *CarController) list(c echo.Context) error {
	user := userClaimsFromContext(c)
	limitStr := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitStr)

	if limitStr == "" {
		limit = 15
	} else if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "limit must be an integer greater than 0")
	}

	lastID := c.QueryParam("last")

	cars, err := cc.service.GetAll(lastID, limit, user.AuthLevel)

	if err != nil {
		if err == services.ErrNotAllowed {
			return echo.NewHTTPError(http.StatusForbidden, err)
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": cars,
	})
}

func (cc *CarController) show(c echo.Context) error {
	id := c.Param("id")

	car, err := cc.service.GetCar(id)

	if err != nil {
		if err == postgres.ErrNoCarFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": car,
	})
}

func (cc *CarController) create(c echo.Context) error {
	var body models.Car

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userClaims := userClaimsFromContext(c)

	car, err := cc.service.AddCar(body, userClaims.ID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"data": car,
	})
}

func (cc *CarController) remove(c echo.Context) error {
	id := c.Param("id")

	userClaims := userClaimsFromContext(c)

	if err := cc.service.DeleteCar(id, userClaims.ID); err != nil {
		if err == services.ErrNotCarOwner {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
