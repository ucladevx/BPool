package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type (
	// UserService is used to handle the user use cases
	UserService interface {
		Login(string) (string, error)
	}

	authCookieInfo struct {
		numDaysValid int
		cookieName   string
	}

	// User http adapter
	User struct {
		logger     Logger
		service    UserService
		authCookie authCookieInfo
	}

	userLoginRequest struct {
		Token string `json:"token"`
	}
)

// NewUserController creates a new auth controller
func NewUserController(u UserService, daysTokenValidFor int, cookieName string, l Logger) *User {
	a := authCookieInfo{daysTokenValidFor, cookieName}

	return &User{
		logger:     l,
		service:    u,
		authCookie: a,
	}
}

// MountRoutes mounts the auth routes
func (u *User) MountRoutes(c *echo.Group) {
	c.POST("/login", u.login)
}

func (u *User) login(c echo.Context) error {
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
