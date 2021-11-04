// Package env holds all environment variable related code
package env_test

import (
	"testing"

	"github.com/rakin92/go-rest-service/pkg/env"
	"github.com/stretchr/testify/assert"
)

func TestMustGet(t *testing.T) {
	t.Run("Panic when can't find env variable", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustGet should have panicked!")
				}
			}()
			// This function should cause a panic
			// should add checks to confirm logger gets called
			env.MustGet("host")
		}()
	})
	t.Run("Returns env variable when found", func(t *testing.T) {
		t.Setenv("host", "localhost")

		got := env.MustGet("host")
		assert.Equal(t, "localhost", got)
	})
}

func TestMustGetBool(t *testing.T) {
	t.Run("Panic when can't find env variable", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustGetBool should have panicked!")
				}
			}()
			// This function should cause a panic
			// should add checks to confirm logger gets called
			env.MustGetBool("debug")
		}()
	})
	t.Run("Panic when fail to parse bool", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustGetBool should have panicked!")
				}
			}()
			// This function should cause a panic
			// should add checks to confirm logger gets called
			t.Setenv("debug", "not_supported_bool")
			env.MustGetBool("debug")
		}()
	})
	t.Run("Returns a valid boolean", func(t *testing.T) {
		type args struct {
			key   string
			value string
		}
		tests := []struct {
			name string
			args args
			want bool
		}{
			{
				name: "TRUE",
				args: args{key: "debug", value: "TRUE"},
				want: true,
			},
			{
				name: "True",
				args: args{key: "debug", value: "True"},
				want: true,
			},
			{
				name: "T",
				args: args{key: "debug", value: "T"},
				want: true,
			},
			{
				name: "t",
				args: args{key: "debug", value: "t"},
				want: true,
			},
			{
				name: "1",
				args: args{key: "debug", value: "1"},
				want: true,
			},
			{
				name: "FALSE",
				args: args{key: "debug", value: "FALSE"},
				want: false,
			},
			{
				name: "False",
				args: args{key: "debug", value: "False"},
				want: false,
			},
			{
				name: "F",
				args: args{key: "debug", value: "F"},
				want: false,
			},
			{
				name: "f",
				args: args{key: "debug", value: "f"},
				want: false,
			},
			{
				name: "0",
				args: args{key: "debug", value: "0"},
				want: false,
			},
		}
		for _, tt := range tests {
			t.Setenv(tt.args.key, tt.args.value)
			t.Run(tt.name, func(t *testing.T) {
				if got := env.MustGetBool(tt.args.key); got != tt.want {
					t.Errorf("MustGetBool() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestMustGetInt32(t *testing.T) {
	t.Run("Panic when can't find env variable", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustGetInt32 should have panicked!")
				}
			}()
			// This function should cause a panic
			// should add checks to confirm logger gets called
			env.MustGetInt32("port")
		}()
	})
	t.Run("Panic when fail to parse int", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustGetInt32 should have panicked!")
				}
			}()
			// This function should cause a panic
			// should add checks to confirm logger gets called
			t.Setenv("port", "not_an_int")
			env.MustGetInt32("port")
		}()
	})
	t.Run("Returns valid int", func(t *testing.T) {
		t.Setenv("port", "8000")

		got := env.MustGetInt32("port")
		assert.Equal(t, 8000, got)
	})
}

func TestMustGetInt64(t *testing.T) {
	t.Run("Panic when can't find env variable", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustGetInt64 should have panicked!")
				}
			}()
			// This function should cause a panic
			// should add checks to confirm logger gets called
			env.MustGetInt64("port")
		}()
	})
	t.Run("Panic when fail to parse int", func(t *testing.T) {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustGetInt64 should have panicked!")
				}
			}()
			// This function should cause a panic
			// should add checks to confirm logger gets called
			t.Setenv("port", "not_an_int")
			env.MustGetInt64("port")
		}()
	})
	t.Run("Returns valid int64", func(t *testing.T) {
		t.Setenv("port", "8000")

		expect := int64(8000)
		got := env.MustGetInt64("port")
		assert.Equal(t, expect, got)
	})
}
