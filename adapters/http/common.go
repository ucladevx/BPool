package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/utils/auth"
)

var (
	// ErrNotAllowed occurs when the user does not have sufficient auth level
	ErrNotAllowed = echo.NewHTTPError(http.StatusForbidden, "user not allowed")
)

func userClaimsFromContext(c echo.Context) *auth.UserClaims {
	possibleUser := c.Get("user")

	user, ok := possibleUser.(auth.UserClaims)
	if !ok {
		return nil
	}

	return &user
}
