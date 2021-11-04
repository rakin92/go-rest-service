// Package consts all constants
package consts

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "convert BaseToSnakeCase",
			args: args{str: "BaseToSnakeCase"},
			want: "base_to_snake_case",
		},
		{
			name: "convert TwoWord",
			args: args{str: "TwoWord"},
			want: "two_word",
		},
		{
			name: "convert smallAndCapital",
			args: args{str: "smallAndCapital"},
			want: "small_and_capital",
		},
		{
			name: "convert already_snake_case",
			args: args{str: "already_snake_case"},
			want: "already_snake_case",
		},
		{
			name: "convert empty",
			args: args{str: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSnakeCase(tt.args.str); got != tt.want {
				t.Errorf("ToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatPermissionDesc(t *testing.T) {
	type args struct {
		action string
		entity string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "create user permission description",
			args: args{action: "create", entity: "users"},
			want: "Allows the user to create users",
		},
		{
			name: "empty param permission description",
			args: args{action: "", entity: ""},
			want: "Allows the user to  ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatPermissionDesc(tt.args.action, tt.args.entity); got != tt.want {
				t.Errorf("FormatPermissionDesc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatPermissionTag(t *testing.T) {
	type args struct {
		action string
		entity string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "create user permission tag",
			args: args{action: "create", entity: "users"},
			want: "create:users",
		},
		{
			name: "update apiKey permission tag",
			args: args{action: "update", entity: "apiKey"},
			want: "update:apiKey",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatPermissionTag(tt.args.action, tt.args.entity); got != tt.want {
				t.Errorf("FormatPermissionTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTableName(t *testing.T) {
	type args struct {
		tablename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "create user permission tag",
			args: args{tablename: "UsersProfile"},
			want: "users_profile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTableName(tt.args.tablename); got != tt.want {
				t.Errorf("GetTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
