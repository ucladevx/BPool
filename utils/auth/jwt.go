package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrTokenInvalid is an invalid token error
	ErrTokenInvalid = errors.New("invalid token")
	// ErrTokenWrongIssuer occurs when token has wrong issuer
	ErrTokenWrongIssuer = errors.New("token has wrong issuer")
)

type (
	// Logger application logger
	Logger interface {
		Debug(args ...interface{})
		Error(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Panic(args ...interface{})
	}

	//Tokenizer creates/parses JWT tokens
	Tokenizer struct {
		secret    []byte
		cookie    string
		issuer    string
		daysValid int
		parser    *jwt.Parser
		logger    Logger
	}
)

// NewTokenizer creates a tokenizer to create/parse JWT tokens
func NewTokenizer(secret, issuer string, daysValid int, l Logger) *Tokenizer {
	return &Tokenizer{
		secret:    []byte(secret),
		issuer:    issuer,
		daysValid: daysValid,
		parser:    &jwt.Parser{},
		logger:    l,
	}
}

// NewJWTmiddleware is an auth middleware to check JWT tokens
func NewJWTmiddleware(t *Tokenizer, cookie string, l Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := c.Cookie(cookie)
			if err != nil {
				l.Info("JWT Middleware - no token", "error", err.Error())
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			claims, err := t.Validate(token.Value)
			if err != nil {
				l.Info("JWT Middleware", "info", "No token")

				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			c.Set("user", claims)

			return next(c)
		}
	}
}

// NewToken creates a new token with the given claims
func (t *Tokenizer) NewToken(claims map[string]interface{}) (string, error) {
	jwtClaims := jwt.MapClaims{}
	jwtClaims["exp"] = time.Now().AddDate(0, 0, t.daysValid).Unix()
	jwtClaims["iss"] = t.issuer

	for key, val := range claims {
		jwtClaims[key] = val
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := token.SignedString(t.secret)
	if err != nil {
		t.logger.Error("Tokenizer.NewToken - sign token", "error", err.Error())
		return "", err
	}

	return signedToken, nil
}

// Validate checks if a token is valid and returns the token's claims
func (t *Tokenizer) Validate(tokenString string) (map[string]interface{}, error) {
	token, err := t.parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.logger.Error("Tokenizer.Validate - token parse",
				"error", "token signing method invalid",
			)

			return nil, ErrTokenInvalid
		}
		return t.secret, nil
	})

	if err != nil {
		t.logger.Error("Tokenizer.Validate - checking token",
			"info", "Token was invalid",
			"error", err.Error(),
		)

		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims.VerifyIssuer(t.issuer, true) {
			return claims, nil
		}

		t.logger.Error("Tokenizer.Validate - checking token",
			"error", "invalid issuer",
		)

		return nil, ErrTokenWrongIssuer
	}

	t.logger.Error("Tokenizer.Validate",
		"info", "Token was invalid",
	)

	return nil, ErrTokenInvalid
}
