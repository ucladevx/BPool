package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"

	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/utils/auth"
)

func TestNewUserLogin(t *testing.T) {
	mockLogger := mocks.Logger{}
	authorizer := auth.NewGoogleAuthorizer("123", "123", "localhost:3000", "test", mockLogger)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/login", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	authorizer.NewUserLogin(c)

	if res.Code != http.StatusFound {
		t.Errorf("expected status code to be 302 found, instead %d", res.Code)
	}

	if len(res.HeaderMap["Set-Cookie"]) == 0 {
		t.Error("should have set a cookie on the response")
	}
}
