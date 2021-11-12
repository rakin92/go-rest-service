package models_test

import (
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/rakin92/go-rest-service/internal/orm/models"
)

func TestGothUserToDBUserProfile(t *testing.T) {
	testID := "ID"
	type args struct {
		i      *goth.User
		update bool
		ids    []int
	}
	tests := []struct {
		name    string
		args    args
		wantO   *models.UserProfile
		wantErr bool
	}{
		{
			name: "missing email create",
			args: args{
				i: &goth.User{
					UserID: testID,
				},
				update: false,
				ids:    []int{},
			},
			wantO:   nil,
			wantErr: true,
		},
		{
			name: "missing id create",
			args: args{
				i: &goth.User{
					Email: "email",
				},
				update: false,
				ids:    []int{},
			},
			wantO:   nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				i: &goth.User{
					UserID: testID,
					Email:  "email",
				},
				update: false,
				ids:    []int{},
			},
			wantO: &models.UserProfile{
				Email:          "email",
				ExternalUserID: testID,
			},
			wantErr: false,
		},
		{
			name: "success with id",
			args: args{
				i: &goth.User{
					UserID: testID,
					Email:  "email",
				},
				update: false,
				ids:    []int{1},
			},
			wantO: &models.UserProfile{
				Email:          "email",
				ExternalUserID: testID,
				BaseModelSeq: models.BaseModelSeq{
					ID: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, err := models.GothUserToDBUserProfile(tt.args.i, tt.args.update, tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GothUserToDBUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("GothUserToDBUserProfile() = %v, want %v", gotO, tt.wantO)
			}
		})
	}
}

func TestGothUserToDBUser(t *testing.T) {
	testID := "ID"
	firstName := "fname"
	lastName := "lname"
	uid, _ := uuid.NewV4()
	type args struct {
		i      *goth.User
		update bool
		ids    []string
	}
	tests := []struct {
		name    string
		args    args
		wantO   *models.User
		wantErr bool
	}{
		{
			name: "missing email create",
			args: args{
				i: &goth.User{
					UserID: testID,
				},
				update: false,
				ids:    []string{},
			},
			wantO:   nil,
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				i: &goth.User{
					Email:     "email",
					FirstName: firstName,
					LastName:  lastName,
				},
				update: false,
				ids:    []string{},
			},
			wantO: &models.User{
				Email:     "email",
				FirstName: &firstName,
				LastName:  &lastName,
			},
			wantErr: false,
		},
		{
			name: "success with id",
			args: args{
				i: &goth.User{
					Email:     "email",
					FirstName: firstName,
					LastName:  lastName,
				},
				update: false,
				ids:    []string{uid.String()},
			},
			wantO: &models.User{
				Email:     "email",
				FirstName: &firstName,
				LastName:  &lastName,
				BaseModelSoftDelete: models.BaseModelSoftDelete{
					BaseModel: models.BaseModel{
						ID: uid,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "fail with id",
			args: args{
				i: &goth.User{
					Email:     "email",
					FirstName: firstName,
					LastName:  lastName,
				},
				update: false,
				ids:    []string{"bad_uuid"},
			},
			wantO:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, err := models.GothUserToDBUser(tt.args.i, tt.args.update, tt.args.ids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GothUserToDBUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotO, tt.wantO) {
				t.Errorf("GothUserToDBUser() = %v, want %v", gotO, tt.wantO)
			}
		})
	}
}
