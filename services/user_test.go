package services_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"gopkg.in/h2non/gock.v1"

	"github.com/stretchr/testify/assert"

	"github.com/ucladevx/BPool/mocks"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/services"
	"github.com/ucladevx/BPool/stores/postgres"
	"github.com/ucladevx/BPool/utils/auth"
)

var (
	johnDoe = models.User{
		ID:           "abc123",
		FirstName:    "JOHN",
		LastName:     "DOE",
		Email:        "johndoe@g.ucla.edu",
		ProfileImage: "ucladevx.com",
		AuthLevel:    services.UserLevel,
	}

	janeSmith = models.User{
		ID:           "123abc",
		FirstName:    "JANE",
		LastName:     "SMITH",
		Email:        "janesmith@g.ucla.edu",
		ProfileImage: "google.com",
		AuthLevel:    services.AdminLevel,
	}
)

func newUserService(store *mocks.UserStore) *services.UserService {
	logger := mocks.Logger{}
	tokenizer := auth.NewTokenizer(
		"secret",
		"bpool",
		14,
		logger,
	)

	return services.NewUserService(store, tokenizer, logger)
}

func TestUserLogin(t *testing.T) {
	store := new(mocks.UserStore)
	service := newUserService(store)
	assert := assert.New(t)

	token := "abcdefg"

	gock.New("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token).
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"hd":             "g.ucla.edu",
			"email":          johnDoe.Email,
			"email_verified": "true",
			"exp":            1522908185,
			"name":           "John Doe",
			"picture":        "ucladevx.com",
			"given_name":     "JOHN",
			"family_name":    "DOE",
		})

	defer gock.Off()

	// first case, assume user in DB
	store.On("GetByEmail", johnDoe.Email).Return(&johnDoe, nil)

	authToken, err := service.Login(token)

	assert.Nil(err, "there should be no error for a valid login and found in DB")
	assert.NotEqual("", authToken, "there should be an auth token for valid user login and found in DB")
	store.AssertExpectations(t)

	secondToken := "1234abcd"

	gock.New("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + secondToken).
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"hd":             "g.ucla.edu",
			"email":          janeSmith.Email,
			"email_verified": "true",
			"exp":            1522908185,
			"name":           "Jane Smith",
			"picture":        "ucladevx.com",
			"given_name":     janeSmith.FirstName,
			"family_name":    janeSmith.LastName,
		})

	store.On("GetByEmail", janeSmith.Email).Return(nil, nil)
	store.On("Insert", mock.AnythingOfType("*models.User")).Return(nil).
		Run(func(args mock.Arguments) {
			// adjust the pointer passed by insert
			arg := args.Get(0).(*models.User)
			*arg = janeSmith
		})

	authToken, err = service.Login(secondToken)

	assert.Nil(err, "there should be no error for a valid login when not found in DB")
	assert.NotEqual("", authToken, "there should be an auth token for valid user login and not found in DB")
	store.AssertExpectations(t)
}

func TestUserGet(t *testing.T) {
	store := new(mocks.UserStore)
	service := newUserService(store)
	assert := assert.New(t)

	store.On("GetByID", "abc").Return(nil, postgres.ErrNoUserFound)

	noUser, err := service.Get("abc")

	assert.Nil(noUser, "for a bad id there should be no user")
	assert.Equal(postgres.ErrNoUserFound, err, "if no user found should return no user found error")

	u := johnDoe
	validID := "test1234"
	u.ID = validID

	store.On("GetByID", validID).Return(&u, nil)

	user, err := service.Get(validID)
	assert.NotNil(user, "user should not be nil for a valid id")
	assert.Nil(err, "for a valid user, error should be nil")
}

func TestUserGetAll(t *testing.T) {
	store := new(mocks.UserStore)
	service := newUserService(store)
	assert := assert.New(t)

	noUsers, err := service.GetAll("", 15, services.UserLevel)
	assert.Nil(noUsers, "when a user does not have the right auth level there should be no users")
	assert.Equal(services.ErrNotAllowed, err, "when a user does not have the right auth level there should be a not allowed error")

	badLimit := -1
	store.On("GetAll", "", 15).Return([]*models.User{&johnDoe, &janeSmith}, nil)

	users, err := service.GetAll("", badLimit, services.AdminLevel)

	assert.Nil(err, "for a bad limit, should still return no error")
	assert.Equal(2, len(users), "the returned users should have length 2")

	store.AssertExpectations(t)
}
