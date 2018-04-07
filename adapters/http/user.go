package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
)

type (
	// UserService is used to handle the user use cases
	UserService interface {
		Login(googleToken string) (string, error)
		Get(id string) (*models.User, error)
	}

	authCookieInfo struct {
		numDaysValid int
		cookieName   string
	}

	// UserController http adapter
	UserController struct {
		logger     interfaces.Logger
		service    UserService
		authCookie authCookieInfo
	}

	userLoginRequest struct {
		Token string `json:"token"`
	}
)

// NewUserController creates a new auth controller
func NewUserController(u UserService, daysTokenValidFor int, cookieName string, l interfaces.Logger) *UserController {
	a := authCookieInfo{daysTokenValidFor, cookieName}

	return &UserController{
		logger:     l,
		service:    u,
		authCookie: a,
	}
}

// MountRoutes mounts the auth routes
func (u *UserController) MountRoutes(c *echo.Group) {
	c.GET("/users/:id", u.show)
	c.POST("/login", u.login)
}

func (u *UserController) login(c echo.Context) error {
	var data userLoginRequest
	if err := c.Bind(&data); err != nil {
		echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	token, err := u.service.Login(data.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cookie := new(http.Cookie)

	daysValid := time.Hour * time.Duration((u.authCookie.numDaysValid * 24))
	cookie.Name = u.authCookie.cookieName
	cookie.Value = token
	cookie.Expires = time.Now().Add(daysValid)
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"data": token,
	})
}

func (u *UserController) show(c echo.Context) error {
	id := c.Param("id")

	if id == "@me" {
		userClaims := userClaimsFromContext(c)
		if userClaims == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "user is not logged in")
		}

		id = userClaims.ID
	}

	user, err := u.service.Get(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": user,
	})
}
