package http

import (
	"github.com/labstack/echo"
	"github.com/ucladevx/BPool/utils/auth"
)

func userClaimsFromContext(c echo.Context) *auth.UserClaims {
	possibleUser := c.Get("user")

	user, ok := possibleUser.(auth.UserClaims)
	if !ok {
		return nil
	}

	return &user
}
