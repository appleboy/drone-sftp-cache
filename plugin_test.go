package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin_check(t *testing.T) {
	tests := []struct {
		name    string
		fields  Plugin
		wantErr bool
	}{
		{
			"Test missing server or username",
			Plugin{
				Server: "localhost",
			},
			true,
		},
		{
			"Test missing password or key",
			Plugin{
				Server:   "localhost",
				Username: "ubuntu",
			},
			true,
		},
	}
	for _, tt := range tests {
		p := &Plugin{
			Server:   tt.fields.Server,
			Username: tt.fields.Username,
			Password: tt.fields.Password,
			Key:      tt.fields.Key,
		}
		if err := p.check(); (err != nil) != tt.wantErr {
			t.Errorf("%q. Plugin.check() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestConfigCheck(t *testing.T) {
	plugin := Plugin{
		Server:   "localhost",
		Username: "drone-scp",
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestIncorrectPassword(t *testing.T) {
	plugin := Plugin{
		Server:   "localhost",
		Username: "drone-scp",
		Port:     "22",
		Password: "123456",
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}
