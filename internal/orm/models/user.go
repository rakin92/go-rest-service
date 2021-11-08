package models

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/rakin92/go-rest-service/pkg/consts"
	"gorm.io/gorm"
)

// ## Entity definitions

// User defines a user for the service
type User struct {
	BaseModelSoftDelete        // We don't to actually delete the users, audit
	Email               string `gorm:"not null;unique;index"`
	FirstName           *string
	LastName            *string
	UserProfiles        []UserProfile `gorm:"association_autocreate:false;association_autoupdate:false"`
	Roles               []Role        `gorm:"many2many:user_roles;association_autocreate:false;association_autoupdate:false"`
	Permissions         []Permission  `gorm:"many2many:user_permissions;association_autocreate:false;association_autoupdate:false"`
}

// UserProfile saves all the related OAuth Profiles
type UserProfile struct {
	BaseModelSeq
	Email          string    `gorm:"uniqueIndex:idx_email_provider_external_user_id"`
	UserID         uuid.UUID `gorm:"not null;index"`
	User           User      `gorm:"association_autocreate:false;association_autoupdate:false"`
	Provider       string    `gorm:"not null;index;unique_index:idx_email_provider_external_user_id;default:'DB'"` // DB means database or no ExternalUserID
	ExternalUserID string    `gorm:"not null;index;unique_index:idx_email_provider_external_user_id"`              // User ID
	FirstName      string
	LastName       string
	AvatarURL      string `gorm:"size:1024"`
	Description    string `gorm:"size:1024"`
}

// UserAPIKey generated api keys for the users
type UserAPIKey struct {
	BaseModelSeqSoftDelete
	Name        string
	User        User         `gorm:"association_autocreate:false;association_autoupdate:false"`
	UserID      uuid.UUID    `gorm:"not null;index"`
	APIKey      string       `gorm:"size:128;unique;index;default:uuid_generate_v4()"`
	Permissions []Permission `gorm:"many2many:user_api_key_permissions;association_autocreate:false;association_autoupdate:false"`
}

// UserRole relation between an user and its roles
type UserRole struct {
	UserID uuid.UUID `gorm:"index"`
	RoleID int       `gorm:"index"`
}

// UserPermission relation between an user and its permissions
type UserPermission struct {
	UserID       uuid.UUID `gorm:"index"`
	PermissionID int       `gorm:"index"`
}

// ## Hooks

// BeforeSave hook for User
func (u *User) BeforeSave(db *gorm.DB) error {
	if u.Email != "" {
		u.Email = strings.ToLower(u.Email)
	}
	return nil
}

// BeforeSave hook for UserAPIKey
func (k *UserAPIKey) BeforeSave(db *gorm.DB) error {
	if k.Name == "" {
		u := &User{}
		if err := db.Where("id = ?", k.UserID).First(u).Error; err != nil {
			return err
		}
	}
	return nil
}

// ## Helper functions

// HasRole verifies if user possesses a role
func (u *User) HasRole(roleID int) (bool, error) {
	for _, r := range u.Roles {
		if r.ID == uint(roleID) {
			return true, nil
		}
	}
	return false, fmt.Errorf("user has no [%d] roleID", roleID)
}

// HasPermission verifies if user has a specific permission
func (u *User) HasPermission(permission string, entity string) (bool, error) {
	tag := fmt.Sprintf("%s:%s", permission, consts.GetTableName(entity))
	for _, p := range u.Permissions {
		if p.Tag == tag {
			return true, nil
		}
	}
	return false, fmt.Errorf("user has no permission: [%s]", tag)
}

// HasPermissionTag verifies if user has a specific permission tag
func (u *User) HasPermissionTag(tag string) (bool, error) {
	for _, r := range u.Permissions {
		if r.Tag == tag {
			return true, nil
		}
	}
	return false, fmt.Errorf("user has no [%s] permission", tag)
}

// CanUpdate verifies if user can update if owner - returns t/f
func (u *User) CanUpdate(id string) (bool, error) {
	if id == u.ID.String() {
		return true, nil
	}
	return false, fmt.Errorf("user [%s] is not the owner", id)
}
