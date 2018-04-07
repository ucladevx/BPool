package auth

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/interfaces"
)

// UserClaims are the claims stored in the JWT
type UserClaims struct {
	ID        string
	Email     string
	AuthLevel int
}

// NewAuthmiddleWare enforces auth on routes
func NewAuthmiddleWare(LowestAuthLevelRequired int, l interfaces.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := userFromContext(c)
			if user == nil {
				l.Error("Auth Middleware - not allowed", "error", "user not logged in")
				return echo.NewHTTPError(http.StatusUnauthorized, "user not logged in")
			}

			if user.AuthLevel < LowestAuthLevelRequired {
				l.Error("Auth Middleware - not authorized", "user_level", user.AuthLevel)
				return echo.NewHTTPError(http.StatusUnauthorized, "you are not authorized")
			}

			return next(c)
		}
	}
}

func userFromContext(c echo.Context) *UserClaims {
	potentialUser := c.Get("user")

	u, ok := potentialUser.(UserClaims)

	if ok {
		return &u
	}

	return nil
}
