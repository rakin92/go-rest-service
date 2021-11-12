package models_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/rakin92/go-rest-service/internal/orm/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func mockOrm(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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
	return gormDB, mock
}

func TestUser_HasPermissionTag(t *testing.T) {
	type fields struct {
		Permissions []models.Permission
	}
	type args struct {
		tag string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "User has permission tag",
			fields: fields{
				Permissions: []models.Permission{
					{Tag: "create:user"},
				},
			},
			args: args{
				tag: "create:user",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "User has permission tag",
			fields: fields{
				Permissions: []models.Permission{
					{Tag: "create:user"},
				},
			},
			args: args{
				tag: "delete:user",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &models.User{
				Permissions: tt.fields.Permissions,
			}
			got, err := u.HasPermissionTag(tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.HasPermissionTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("User.HasPermissionTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	type fields struct {
		Roles []models.Role
	}
	type args struct {
		roleID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "User has role",
			fields: fields{
				Roles: []models.Role{
					{BaseModelSeq: models.BaseModelSeq{ID: 1}},
				},
			},
			args:    args{roleID: 1},
			want:    true,
			wantErr: false,
		},
		{
			name: "User has role",
			fields: fields{
				Roles: []models.Role{
					{BaseModelSeq: models.BaseModelSeq{ID: 1}},
				},
			},
			args:    args{roleID: 2},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &models.User{
				Roles: tt.fields.Roles,
			}
			got, err := u.HasRole(tt.args.roleID)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.HasRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("User.HasRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_CanUpdate(t *testing.T) {
	u, err := uuid.NewV4()
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		BaseModelSoftDelete models.BaseModelSoftDelete
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "User can update",
			fields: fields{
				BaseModelSoftDelete: models.BaseModelSoftDelete{
					BaseModel: models.BaseModel{ID: u},
				},
			},
			args:    args{id: u.String()},
			want:    true,
			wantErr: false,
		},
		{
			name: "User can't update",
			fields: fields{
				BaseModelSoftDelete: models.BaseModelSoftDelete{
					BaseModel: models.BaseModel{ID: u},
				},
			},
			args:    args{id: "not_user_uuid"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &models.User{
				BaseModelSoftDelete: tt.fields.BaseModelSoftDelete,
			}
			got, err := u.CanUpdate(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.CanUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("User.CanUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	type fields struct {
		Permissions []models.Permission
	}
	type args struct {
		permission string
		entity     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "User has permission",
			fields: fields{
				Permissions: []models.Permission{
					{Tag: "create:user"},
				},
			},
			args:    args{permission: "create", entity: "user"},
			want:    true,
			wantErr: false,
		},
		{
			name: "User has permission",
			fields: fields{
				Permissions: []models.Permission{
					{Tag: "create:user"},
				},
			},
			args:    args{permission: "delete", entity: "user"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &models.User{
				Permissions: tt.fields.Permissions,
			}
			got, err := u.HasPermission(tt.args.permission, tt.args.entity)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.HasPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("User.HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_BeforeSave(t *testing.T) {
	gormDB, _ := mockOrm(t)
	type fields struct {
		Email string
	}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantEmail string
	}{
		{
			name:      "uppercase to lowercase email",
			fields:    fields{Email: "SOME_EMAIL@EMAIL.COM"},
			args:      args{db: gormDB},
			wantErr:   false,
			wantEmail: "some_email@email.com",
		},
		{
			name:      "lowercase to lowercase email",
			fields:    fields{Email: "SOME_EMAIL@EMAIL.COM"},
			args:      args{db: gormDB},
			wantErr:   false,
			wantEmail: "some_email@email.com",
		},
		{
			name:      "camel to lowercase email",
			fields:    fields{Email: "SomeEmail@EMAIL.COM"},
			args:      args{db: gormDB},
			wantErr:   false,
			wantEmail: "someemail@email.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &models.User{
				Email: tt.fields.Email,
			}
			if err := u.BeforeSave(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("User.BeforeSave() error = %v, wantErr %v", err, tt.wantErr)
			}
			if u.Email != tt.wantEmail {
				t.Errorf("User.BeforeSave() email = %v, wantEmail %v", u.Email, tt.wantErr)
			}
		})
	}
}
