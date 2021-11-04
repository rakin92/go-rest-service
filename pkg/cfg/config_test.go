// Package cfg is the configuration package hold all config objects
package cfg

import (
	"testing"
)

func Test_getValidHost(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default localhost with (:)",
			args: args{host: ":"},
			want: "localhost",
		},
		{
			name: "non default host",
			args: args{host: "foo.bar.com"},
			want: "foo.bar.com",
		},
		{
			name: "empty host",
			args: args{host: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getValidHost(tt.args.host); got != tt.want {
				t.Errorf("getValidHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ListenEndpoint(t *testing.T) {
	type fields struct {
		Host string
		Port string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "prod port 80",
			fields: fields{
				Port: "80",
				Host: "foo.bar.com",
			},
			want: "foo.bar.com",
		},
		{
			name: "dev port 7777",
			fields: fields{
				Port: "7777",
				Host: "localhost",
			},
			want: "localhost:7777",
		},
		{
			name: "default : localhost",
			fields: fields{
				Port: "7777",
				Host: ":",
			},
			want: ":7777",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if got := s.ListenEndpoint(); got != tt.want {
				t.Errorf("Server.ListenEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_VersionedEndpoint(t *testing.T) {
	type fields struct {
		ServiceVersion string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "with version v1",
			fields: fields{
				ServiceVersion: "v1",
			},
			args: args{path: "/foo"},
			want: "/v1/foo",
		},
		{
			name: "with version v2",
			fields: fields{
				ServiceVersion: "v1",
			},
			args: args{path: "/foo"},
			want: "/v1/foo",
		},
		{
			name: "with version v1",
			fields: fields{
				ServiceVersion: "v1",
			},
			args: args{path: "/foo"},
			want: "/v1/foo",
		},
		{
			name:   "missing version default v1",
			fields: fields{},
			args:   args{path: "/foo"},
			want:   "/v1/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				ServiceVersion: tt.fields.ServiceVersion,
			}
			if got := s.VersionedEndpoint(tt.args.path); got != tt.want {
				t.Errorf("Server.VersionedEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_SchemaVersionedEndpoint(t *testing.T) {
	type fields struct {
		Host           string
		Port           string
		URISchema      string
		ServiceVersion string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "port 80",
			fields: fields{
				Port:           "80",
				Host:           "foo.bar.com",
				URISchema:      "https://",
				ServiceVersion: "v1",
			},
			args: args{path: "/baz"},
			want: "https://foo.bar.com/v1/baz",
		},
		{
			name: "port dev port 7777",
			fields: fields{
				Port:           "7777",
				Host:           "localhost",
				URISchema:      "http://",
				ServiceVersion: "v1",
			},
			args: args{path: "/baz"},
			want: "http://localhost:7777/v1/baz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Host:           tt.fields.Host,
				Port:           tt.fields.Port,
				URISchema:      tt.fields.URISchema,
				ServiceVersion: tt.fields.ServiceVersion,
			}
			if got := s.SchemaVersionedEndpoint(tt.args.path); got != tt.want {
				t.Errorf("Server.SchemaVersionedEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
