package models

import (
	"reflect"
	"testing"

	"github.com/markbates/goth"
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
		wantO   *UserProfile
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
			wantO: &UserProfile{
				Email:          "email",
				ExternalUserID: testID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, err := GothUserToDBUserProfile(tt.args.i, tt.args.update, tt.args.ids...)
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
	type args struct {
		i      *goth.User
		update bool
		ids    []string
	}
	tests := []struct {
		name    string
		args    args
		wantO   *User
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
			wantO: &User{
				Email:     "email",
				FirstName: &firstName,
				LastName:  &lastName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotO, err := GothUserToDBUser(tt.args.i, tt.args.update, tt.args.ids...)
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
