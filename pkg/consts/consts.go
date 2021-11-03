// Package consts all constants
package consts

import (
	"fmt"
	"regexp"
	"strings"
)

type permissionTypes struct {
	Create string
	Read   string
	Update string
	Delete string
	List   string
	Assign string
	Upload string
}

type entitynames struct {
	Users           string
	Roles           string
	Permissions     string
	RoleParents     string
	RolePermissions string
	UserPermissions string
	UserProfiles    string
	UserRoles       string
}

type role struct {
	Name        string
	Description string
}

type dialects struct {
	PostgresSQL string
	MySQL       string
}

var (
	// Permissions has the types of permissions that can be assigned
	Permissions = permissionTypes{
		Create: "create:%s",
		Read:   "read:%s",
		Update: "update:%s",
		Delete: "delete:%s",
		List:   "list:%s",
		Assign: "assign:%s",
		Upload: "upload:%s",
	}

	// EntityNames the names of the tables in the server
	EntityNames = entitynames{
		Users:           "Users",
		Roles:           "Roles",
		Permissions:     "Permissions",
		RoleParents:     "RoleParents",
		RolePermissions: "RolePermissions",
		UserPermissions: "UserPermissions",
		UserProfiles:    "UserProfiles",
		UserRoles:       "UserRoles",
	}
	// Dialects are definition of databases
	Dialects = dialects{
		PostgresSQL: "postgres",
		MySQL:       "mysql",
	}

	// Roles that are part of the systme
	Roles = []role{
		{
			Name:        "admin",
			Description: "Administrator of the app",
		},
		{
			Name:        "user",
			Description: "Normal user of the app",
		},
	}

	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// GetTableName gets the db normalized tablename
func GetTableName(tablename string) string {
	return ToSnakeCase(tablename)
}

// FormatPermissionTag returns a string formatted action:entity permission
func FormatPermissionTag(action string, entity string) string {
	return fmt.Sprintf(action, entity)
}

// FormatPermissionDesc returns a string with the description of the
// action:entity permission
func FormatPermissionDesc(action string, entity string) string {
	return "Allows the user to " +
		strings.ReplaceAll(FormatPermissionTag(action, entity), ":", " ")
}

// ToSnakeCase converts camelcase str to snake_case
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ContextKey defines a type for context keys shared in the app
type ContextKey string

// ContextKeys holds the context keys throughout the project
type ContextKeys struct {
	GothicProviderCtxKey ContextKey // Provider for Gothic library
	ProviderCtxKey       ContextKey // Provider in Auth
	UserCtxKey           ContextKey // User db object in Auth
}

var (
	// ProjectContextKeys the project's context keys
	ProjectContextKeys = ContextKeys{
		GothicProviderCtxKey: "provider",
		ProviderCtxKey:       "gg-provider",
		UserCtxKey:           "gg-auth-user",
	}
)
