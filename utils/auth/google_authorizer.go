package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleInfoEndpoint = "https://www.googleapis.com/oauth2/v3/userinfo"
	uclaDomain         = "g.ucla.edu"
	tokenSub           = "google oauth"
)

var (
	// ErrGoogleInvalidState is returned when state token is invalid
	ErrGoogleInvalidState = errors.New("invalid state token")
)

type (
	// GoogleUser is the information provided from google
	GoogleUser struct {
		Sub          string `json:"sub"`
		FullName     string `json:"name"`
		FirstName    string `json:"given_name"`
		LastName     string `json:"family_name"`
		Picture      string `json:"picture"`
		Email        string `json:"email"`
		EmailVerifed bool   `json:"email_verified"`
	}

	// GoogleAuthorizer allows for google oAuth
	GoogleAuthorizer struct {
		oAuthConfig *oauth2.Config
		tokenizer   *Tokenizer
		logger      Logger
	}
)

// NewUserLogin generates a unique url for user login
func (a *GoogleAuthorizer) NewUserLogin(c echo.Context) error {
	state := uuid.NewV4().String()

	expiry := time.Now().Add(10 * time.Minute)

	// add additional info here if we want something like a specific redirect
	stateToken, err := a.tokenizer.NewToken(map[string]interface{}{
		"sub": tokenSub,
		"exp": expiry,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	// TODO: make these changeable
	cookie := &http.Cookie{
		Name:     state,
		Value:    stateToken,
		Expires:  expiry,
		HttpOnly: true,
	}

	c.SetCookie(cookie)

	url := a.oAuthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("hd", uclaDomain),
	)

	return c.Redirect(http.StatusFound, url)
}

// GetUserFromCode gets user info from google given oauth code
func (a *GoogleAuthorizer) GetUserFromCode(code, state, stateToken string) (*GoogleUser, error) {
	claims, err := a.tokenizer.Validate(stateToken)
	if err != nil {
		return nil, err
	}

	if claims["sub"] != tokenSub {
		return nil, ErrGoogleInvalidState
	}

	token, err := a.oAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	client := a.oAuthConfig.Client(oauth2.NoContext, token)

	resp, err := client.Get(googleInfoEndpoint)
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	var user GoogleUser
	if err = json.Unmarshal(contents, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// NewGoogleAuthorizer creates an authorizer for google oAuth
func NewGoogleAuthorizer(cID, cSecret, redirectURL, secret string, l Logger) *GoogleAuthorizer {
	client := &oauth2.Config{
		ClientID:     cID,
		ClientSecret: cSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	t := NewTokenizer(secret, "bpool", 0, l)

	return &GoogleAuthorizer{
		oAuthConfig: client,
		tokenizer:   t,
		logger:      l,
	}
}
