package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ucladevx/BPool/interfaces"
)

const (
	googleInfoEndpoint = "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token="
	uclaDomain         = "g.ucla.edu"
)

var (
	// ErrGoogleInvalidToken is returned when google token is invalid
	ErrGoogleInvalidToken = errors.New("the token was invalid")
	// ErrGoogleError is returned when there is a problem verifiying the token from google
	ErrGoogleError = errors.New("problem verifying the token from google")
	// ErrUserParseError is returned when google sends a JSON object that is unparsable
	ErrUserParseError = errors.New("could not parse the user")
	// ErrWrongHostedDomain occurs when user's hd is not UCLA
	ErrWrongHostedDomain = errors.New("user is not part of UCLA")
)

type (
	// GoogleUser is the information provided from google
	GoogleUser struct {
		Sub       string `json:"sub"`
		HD        string `json:"hd"`
		Email     string `json:"email"`
		FullName  string `json:"name"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Picture   string `json:"picture"`
	}

	// GoogleAuthorizer allows for google oAuth
	GoogleAuthorizer struct {
		logger interfaces.Logger
	}
)

// NewGoogleAuthorizer creates an authorizer for google oAuth
func NewGoogleAuthorizer(l interfaces.Logger) *GoogleAuthorizer {
	return &GoogleAuthorizer{
		logger: l,
	}
}

// UserLogin is used to verify a token and parse the token from google
func (a *GoogleAuthorizer) UserLogin(token string) (*GoogleUser, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	url := googleInfoEndpoint + token

	resp, err := client.Get(url)
	if err != nil {
		a.logger.Error("GoogleAuthorizer.UserLogin - HTTP error", "error", err.Error())
		return nil, ErrGoogleError
	}

	// marshal the JSON user into a google user
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 200 {
		a.logger.Error("GoogleAuthorizer.UserLogin - response error",
			"error", string(contents),
		)
		return nil, ErrGoogleError
	}

	var user GoogleUser
	if err = json.Unmarshal(contents, &user); err != nil {
		a.logger.Error("GoogleAuthorizer.UserLogin - parse err",
			"error", err.Error(),
			"response", string(contents),
		)

		return nil, ErrUserParseError
	}

	if user.HD != uclaDomain {
		return nil, ErrWrongHostedDomain
	}

	return &user, nil
}
