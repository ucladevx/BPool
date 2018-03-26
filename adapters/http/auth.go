package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/utils/auth"
)

// Auth http adapter
type Auth struct {
	authorizer *auth.GoogleAuthorizer
	logger     Logger
}

func (a *Auth) login(c echo.Context) error {
	return a.authorizer.NewUserLogin(c)
}

func (a *Auth) callback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	stateToken, err := c.Cookie(state)
	if err != nil {
		c.JSON(http.StatusUnauthorized, echo.Map{
			"error": err.Error(),
		})
	}
	user, err := a.authorizer.GetUserFromCode(code, state, stateToken.Value)

	if err != nil {
		c.JSON(http.StatusUnauthorized, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"user": user,
	})
}

// MountRoutes mounts the auth routes
func (a *Auth) MountRoutes(c *echo.Group) {
	c.GET("/login", a.login)
	c.GET("/google/callback", a.callback)
}

// NewAuthController creates a new auth controller
func NewAuthController(a *auth.GoogleAuthorizer, l Logger) *Auth {
	return &Auth{
		authorizer: a,
		logger:     l,
	}
}
