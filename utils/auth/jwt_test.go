package auth_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo"

	"github.com/dgrijalva/jwt-go"
	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/utils/auth"
)

const (
	encodedHeader = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	emptyToken    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.t-IDcSemACt8x4iTMCda8Yhe3iZaWbvV5XKSTbuAn0M"
	withData      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIxNTIxOTk0MjA2Mjg4IiwiaXNzIjoidGVzdCIsInNvbWUiOiJkYXRhIn0.9PPXPwtWKotNw3fPmPa2IrJJizWvrjVLD_UxDD3_Ir4"
)

func TestTokenizerNewToken(t *testing.T) {
	mockLogger := mocks.Logger{}
	tokenizer := auth.NewTokenizer("secret", "test", 1, mockLogger)

	tables := []struct {
		name   string
		claim  map[string]interface{}
		result string
	}{
		{"no claims", nil, emptyToken},
		{"with claims", map[string]interface{}{"some": "data", "exp": "1521994206288"}, withData},
	}

	for _, tt := range tables {
		token, err := tokenizer.NewToken(tt.claim)
		if err != nil {
			t.Errorf("input of %s received error %s", tt.name, tt.result)
		}

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			t.Errorf("input of %s should have length 3, was %d", tt.name, len(parts))
		}

		if parts[0] != encodedHeader {
			t.Errorf("input of %s should have header %s, returned %s", tt.name, encodedHeader, parts[0])
		}

		if tt.claim != nil && tt.result != token {
			t.Errorf("input of %s should have token %s, but token is %s", tt.name, tt.result, token)
		}
	}
}

func TestTokenizerValidate(t *testing.T) {
	mockLogger := mocks.Logger{}
	tokenizer := auth.NewTokenizer("secret", "test", 1, mockLogger)

	tables := []struct {
		name  string
		token string
		claim map[string]interface{}
		err   error
	}{
		{
			"token with no claims",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyNTMxODM5MjU3MjQyIiwiaXNzIjoidGVzdCJ9.2u-gsK6KSMP9pmrCIz44c-VzfPUqxqqWBi26oJEG27Q",
			nil,
			nil,
		},
		{
			"token with claims",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzb21lIjoiZGF0YSIsImV4cCI6IjI1MzE4MzkyNTcyNDIiLCJpc3MiOiJ0ZXN0In0.rX7fAtykYvE2Y_yjHnI1P5YJxN8JYlj5rLixNyJDJ00",
			map[string]interface{}{"some": "data", "exp": "2531839257242"},
			nil,
		},
		{
			"expired token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjIwNDE1ODQsImlzcyI6InRlc3QifQ.fL8T4Yf2ZzNzbEils9WbfRfU6dsRrul_b_qY5ln-Kr4",
			nil,
			errors.New("Token is expired"),
		},
		{
			"wrong issuer token",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyNTMxODM5MjU3MjQyIiwiaXNzIjoibGl0In0.1DBoD15xV62uNbkx2RYxJOFa3TDpzoymlca_ImvhqIw",
			nil,
			auth.ErrTokenWrongIssuer,
		},
		{
			"wrong signature token",
			"eyJhbGciOiJIUzI1NiJ9.e30.WpcxeaF9d7DQP7PeL34f8A6hl5yb1ZOcPD-zYXdx3tw",
			nil,
			jwt.ErrSignatureInvalid,
		},
		{
			"wrong algorithm token",
			"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.e30.Vhe_alDDkRu6VlYIxgT7hWPpKVb98TCa7DoG0XrokdoUmxOf64AGcl6ie9x6XSvCP98oVHtR8yByj1pmp4DS6w",
			nil,
			auth.ErrTokenInvalid,
		},
	}

	for _, tt := range tables {
		claims, err := tokenizer.Validate(tt.token)
		if err != nil {
			if tt.err == nil {
				t.Errorf("%s should not have erred, got error %s", tt.name, err.Error())
			}

			if tt.err != err && tt.err.Error() != err.Error() {
				t.Errorf("%s should have erred with %s, got error %s", tt.name, tt.err.Error(), err.Error())
			}

			continue
		}

		for key, val := range tt.claim {
			inToken := claims[key]
			if inToken != val {
				t.Errorf("%s should have claim %s with val %v, instead token claim is %v", tt.name, key, val, inToken)
			}
		}
	}
}

func TestJWTmiddleware(t *testing.T) {
	mockLogger := mocks.Logger{}
	tokenizer := auth.NewTokenizer("secret", "test", 1, mockLogger)
	e := echo.New()
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "Good")
	}
	f := auth.NewJWTmiddleware(tokenizer, "auth", mockLogger)(h)

	encodedUser := auth.UserClaims{
		ID:        "abcdefg",
		Email:     "test@gmail.com",
		AuthLevel: 0,
	}

	tables := []struct {
		name       string
		token      string
		checkUser  bool
		statusCode int
	}{
		{
			"valid token provided",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyNTMxODM5MjU3MjQyIiwiaXNzIjoidGVzdCIsImlkIjoiYWJjZGVmZyIsImVtYWlsIjoidGVzdEBnbWFpbC5jb20iLCJhdXRoX2xldmVsIjowfQ.Izzlxf3zdPTIRMm2UF0PwC2lLCYUkF3kPwwJwQJMT-U",
			true,
			http.StatusOK,
		},
		{
			"no token provided",
			"",
			false,
			http.StatusOK,
		},
		{
			"expired token provided",
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjIwNDE1ODQsImlzcyI6InRlc3QiLCJpZCI6ImFiY2RlZGZnIiwiZW1haWwiOiJ0ZXN0QGdtYWlsLmNvbSIsImF1dGhfbGV2ZWwiOjB9.rzgIK61wECrtK6V-5fZUOs2-FWE8t4KE2HKG1BbENhk",
			false,
			http.StatusUnauthorized,
		},
	}

	for _, tt := range tables {
		req := httptest.NewRequest(echo.GET, "/login", nil)
		if tt.token != "" {
			c := &http.Cookie{
				Name:  "auth",
				Value: tt.token,
			}

			req.AddCookie(c)
		}

		res := httptest.NewRecorder()
		c := e.NewContext(req, res)

		err := f(c)
		if err != nil {
			he := err.(*echo.HTTPError)
			if he.Code != tt.statusCode {
				t.Errorf("if %s should have responded with %d, instead responded %d", tt.name, tt.statusCode, he.Code)
			}
			continue
		}

		if res.Code != tt.statusCode {
			t.Errorf("if %s should have responded with %d, instead responded %d", tt.name, tt.statusCode, res.Code)
		}

		if tt.checkUser {
			user, ok := c.Get("user").(auth.UserClaims)
			if !ok {
				t.Errorf("%s should have added a user to the context but did not", tt.name)
			}

			if !reflect.DeepEqual(user, encodedUser) {
				t.Errorf("%s should have same user as the encoded user, %v, but was %v", tt.name, encodedUser, user)
			}
		}
	}
}
