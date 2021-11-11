// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/markbates/goth"
	"github.com/rakin92/go-rest-service/internal/orm/models"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestORM_FindUserByAPIKey(t *testing.T) {
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockdb,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer mockdb.Close()

	type fields struct {
		DB   *gorm.DB
		rows *sqlmock.Rows
	}
	type args struct {
		apiKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "valid api key",
			args: args{
				apiKey: "valid",
			},
			fields: fields{
				DB:   gormDB,
				rows: sqlmock.NewRows([]string{"id", "first_name"}).AddRow(1, "name"),
			},
			want:    &models.User{},
			wantErr: false,
		},
		{
			name: "invalid api key",
			args: args{
				apiKey: "invalid",
			},
			fields: fields{
				DB:   gormDB,
				rows: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ORM{
				DB: tt.fields.DB,
			}

			query := mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(tt.args.apiKey)
			if tt.wantErr {
				query.WillReturnError(errors.New("Bad Id"))
			} else {
				query.WillReturnRows(tt.fields.rows)
			}

			got, err := o.FindUserByAPIKey(tt.args.apiKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ORM.FindUserByAPIKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ORM.FindUserByAPIKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestORM_FindUserByJWT(t *testing.T) {
	mockdb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockdb,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer mockdb.Close()

	type fields struct {
		DB   *gorm.DB
		rows *sqlmock.Rows
	}
	type args struct {
		email    string
		provider string
		userID   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "valid user",
			args: args{
				email:    "valid",
				userID:   "valid",
				provider: "valid",
			},
			fields: fields{
				DB:   gormDB,
				rows: sqlmock.NewRows([]string{"id", "first_name"}).AddRow(1, "name"),
			},
			want:    &models.User{},
			wantErr: false,
		},
		{
			name: "valid user sql error",
			args: args{
				email:    "valid",
				userID:   "valid",
				provider: "valid",
			},
			fields: fields{
				DB:   gormDB,
				rows: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing user id",
			args: args{
				email:    "valid",
				userID:   "",
				provider: "valid",
			},
			fields: fields{
				DB:   gormDB,
				rows: sqlmock.NewRows([]string{}).AddRow(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "missing user email",
			args: args{
				email:    "",
				userID:   "valid",
				provider: "valid",
			},
			fields: fields{
				DB:   gormDB,
				rows: sqlmock.NewRows([]string{}).AddRow(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ORM{
				DB: tt.fields.DB,
			}

			mock.ExpectBegin()
			query := mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).WithArgs(tt.args.email, tt.args.provider, tt.args.userID)
			if tt.wantErr {
				query.WillReturnError(errors.New("Bad Id"))
			} else {
				query.WillReturnRows(tt.fields.rows)
			}

			got, err := o.FindUserByJWT(tt.args.email, tt.args.provider, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ORM.FindUserByJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ORM.FindUserByJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestORM_UpsertUserProfile(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		gu *goth.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ORM{
				DB: tt.fields.DB,
			}
			got, err := o.UpsertUserProfile(tt.args.gu)
			if (err != nil) != tt.wantErr {
				t.Errorf("ORM.UpsertUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ORM.UpsertUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		c *cfg.DB
	}
	tests := []struct {
		name    string
		args    args
		want    *ORM
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}
