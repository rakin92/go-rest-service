package conv

import (
	"errors"

	"github.com/markbates/goth"
	"github.com/rakin92/go-rest-service/internal/orm/models"
)

// GothUserToDBUser transforms [user] goth to db model
func GothUserToDBUser(i *goth.User, update bool) (o *models.User, err error) {
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &models.User{
		Email:     i.Email,
		FirstName: &i.FirstName,
		LastName:  &i.LastName,
	}

	return o, err
}

// GothUserToDBUserProfile transforms [user] goth to db model
func GothUserToDBUserProfile(i *goth.User, update bool) (o *models.UserProfile, err error) {
	if i.UserID == "" && !update {
		return nil, errors.New("field [UserID] is required")
	}
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &models.UserProfile{
		ExternalUserID: i.UserID,
		Provider:       i.Provider,
		Email:          i.Email,
		FirstName:      i.FirstName,
		LastName:       i.LastName,
		AvatarURL:      i.AvatarURL,
	}

	return o, err
}
