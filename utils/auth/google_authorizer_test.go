package auth_test

import (
	"reflect"
	"testing"

	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/utils/auth"

	gock "gopkg.in/h2non/gock.v1"
)

func TestUserLogin(t *testing.T) {
	mockLogger := mocks.Logger{}
	authorizer := auth.NewGoogleAuthorizer(mockLogger)

	defer gock.Off()

	token := "adlfkasdlkfjasld13241234"

	expectedUser := auth.GoogleUser{
		Sub:       "ljadfjaksdfl;asd",
		FullName:  "John Doe",
		FirstName: "JOHN",
		LastName:  "DOE",
		Picture:   "ucladevx.com",
		Email:     "johndoe@g.ucla.edu",
		HD:        "g.ucla.edu",
	}

	gock.New("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token).
		Reply(200).
		JSON(map[string]interface{}{
			"azp":            "341675348183-sq8g1sc1bon2d9k3j11ju72mod2t2gl4.apps.googleusercontent.com",
			"aud":            "341675348183-sq8g1sc1bon2d9k3j11ju72mod2t2gl4.apps.googleusercontent.com",
			"sub":            expectedUser.Sub,
			"hd":             expectedUser.HD,
			"email":          expectedUser.Email,
			"email_verified": "true",
			"at_hash":        "1sadfjlsdfajklasdf",
			"exp":            1522908185,
			"iss":            "accounts.google.com",
			"jti":            "558767f9a60dasfsadfasdfasdf",
			"iat":            1522904585,
			"name":           expectedUser.FullName,
			"picture":        expectedUser.Picture,
			"given_name":     expectedUser.FirstName,
			"family_name":    expectedUser.LastName,
			"locale":         "en",
		})

	user, err := authorizer.UserLogin(token)

	if err != nil {
		t.Errorf("With valid response, there should have been no error, got %s", err.Error())
	}

	if !reflect.DeepEqual(expectedUser, *user) {
		t.Error("with valid response, we should have gotten the expected user, we did not.")
	}

	newToken := "asdflkasdfljasdf"

	gock.New("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + newToken).
		Reply(400).
		JSON(map[string]string{
			"error_description": "Invalid Value",
		})
	user, err = authorizer.UserLogin(newToken)
	t.Log(user)
	if err == nil {
		t.Errorf("With invalid response, there should have been an error and no user, but got none")
	}
}
