package services

import (
	"errors"

	"github.com/ucladevx/BPool/interfaces"
	"github.com/ucladevx/BPool/models"
	"github.com/ucladevx/BPool/stores"
	"github.com/ucladevx/BPool/utils/auth"
)

var (
	// ErrDuplciateRating occurs when user tries to give duplicate ratings for a ride
	ErrDuplciateRating = errors.New("The same user cannot give more than one rating for a given ride")

	// ErrNotPassenger occurs when user tries to leave a rating for a ride user did not take
	ErrNotPassenger = errors.New("user is not a passenger of a ride")
)

type (
	// RatingService provides all uses for ratings
	RatingService struct {
		store            RatingStore
		passengerService *PassengerService
		logger           interfaces.Logger
	}

	// RatingStore any store that handles ratings
	RatingStore interface {
		Insert(rating *models.Rating) error
		GetByID(string) (*models.Rating, error)
		WhereMany([]stores.QueryModifier) ([]*models.Rating, error)
		Average(string) (float32, error)
		Delete(string) error
	}
)

// NewRatingService create new RatingService
func NewRatingService(store RatingStore, p *PassengerService, l interfaces.Logger) *RatingService {
	return &RatingService{
		store:            store,
		passengerService: p,
		logger:           l,
	}
}

// Create persists a rating
func (ratingService *RatingService) Create(data models.Rating, user *auth.UserClaims) (*models.Rating, error) {
	rat, rideID, rateeID, raterID, comment := data.Rating, data.RideID, data.RateeID, data.RaterID, data.Comment

	passengers, err := ratingService.passengerService.GetAllByRideID(rideID, user)
	isPassenger := false

	if err != nil {
		return nil, err
	}

	for _, passenger := range passengers {
		if passenger.PassengerID == raterID {
			isPassenger = true
		}
	}

	if isPassenger == false {
		ratingService.logger.Error("RatingService.Create - user must be in ride", "error", err)
		return nil, ErrNotPassenger
	}

	queryModifiers := []stores.QueryModifier{
		stores.QueryMod("ride_id", stores.EQ, rideID),
		stores.QueryMod("rater_id", stores.EQ, raterID),
	}

	duplicates, err := ratingService.store.WhereMany(queryModifiers)

	if err != nil {
		return nil, err
	}

	if len(duplicates) > 0 {
		ratingService.logger.Error("RatingService.Create - cannot add duplicate ratings", "error", err)
		return nil, ErrDuplciateRating
	}

	rating := &models.Rating{
		Rating:  rat,
		RideID:  rideID,
		RateeID: rateeID,
		RaterID: raterID,
		Comment: comment,
	}

	// Model level validation
	if err := rating.Validate(); err != nil {
		ratingService.logger.Error("RatingService.Create - validate", "error", err)
		return nil, err
	}

	// DB level err handling
	if err := ratingService.store.Insert(rating); err != nil {
		ratingService.logger.Error("RatingService.Create - unable to leave rating", "error", err.Error())
		return nil, err
	}

	return rating, nil
}

// GetByID returns average rating for a certain ride
func (ratingService *RatingService) GetByID(id string, user *auth.UserClaims) (*models.Rating, error) {
	if user.AuthLevel != AdminLevel {
		return nil, ErrForbidden
	}

	return ratingService.store.GetByID(id)
}

// GetRatingByUserID returns average rating for a certain ride
func (ratingService *RatingService) GetRatingByUserID(userID string) (float32, error) {
	rating, err := ratingService.store.Average(userID)

	if err != nil {
		ratingService.logger.Error("RatingService.GetRatingByRideID - Average", "error", err.Error())
		return -1, err
	}

	return rating, nil
}

// Delete removes rating from store if user is allowed to
func (ratingService *RatingService) Delete(ratingID string, user *auth.UserClaims) error {
	rating, err := ratingService.store.GetByID(ratingID)

	if err != nil {
		ratingService.logger.Error("RatingService.Delete - Delete", "error", err.Error())
		return err
	}

	if user.AuthLevel != AdminLevel {
		return ErrForbidden
	}

	return ratingService.store.Delete(rating.ID)
}
