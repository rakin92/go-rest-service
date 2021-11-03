package models

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
)

// GothUserToDBUser transforms [user] goth to db model
func GothUserToDBUser(i *goth.User, update bool, ids ...string) (o *User, err error) {
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &User{
		Email:     i.Email,
		FirstName: &i.FirstName,
		LastName:  &i.LastName,
	}
	if len(ids) > 0 {
		updID, err := uuid.FromString(ids[0])
		if err != nil {
			return nil, err
		}
		o.ID = updID
	}
	return o, err
}

// GothUserToDBUserProfile transforms [user] goth to db model
func GothUserToDBUserProfile(i *goth.User, update bool, ids ...int) (o *UserProfile, err error) {
	if i.UserID == "" && !update {
		return nil, errors.New("field [UserID] is required")
	}
	if i.Email == "" && !update {
		return nil, errors.New("field [Email] is required")
	}
	o = &UserProfile{
		ExternalUserID: i.UserID,
		Provider:       i.Provider,
		Email:          i.Email,
		FirstName:      i.FirstName,
		LastName:       i.LastName,
	}
	if len(ids) > 0 {
		updID := ids[0]
		o.ID = updID
	}
	return o, err
}
