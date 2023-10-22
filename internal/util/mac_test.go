package util

import (
	"encoding/hex"
	"reflect"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMAC_Validate(t *testing.T) {
	type args struct {
		message    string
		messageMAC string
	}
	tests := []struct {
		name string
		key  string
		args args
		want bool
	}{
		{
			name: "valid",
			key:  "some key",
			args: args{
				"some text",
				"f65df1c3981b7fd1a9174c5058a5dbf0734a7ffc1357003dbfbed7a623f83fae",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewMAC([]byte(tt.key))
			message := []byte(tt.args.message)
			messageMAC, err := hex.DecodeString(tt.args.messageMAC)
			assert.NilError(t, err)
			if got := h.Validate(message, messageMAC); got != tt.want {
				t.Errorf("MAC.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMAC_Generate(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		message string
		want    string
	}{
		{
			name:    "valid",
			key:     "some key",
			message: "some text",
			want:    "f65df1c3981b7fd1a9174c5058a5dbf0734a7ffc1357003dbfbed7a623f83fae",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewMAC([]byte(tt.key))
			message := []byte(tt.message)
			messageMAC, err := hex.DecodeString(tt.want)
			assert.NilError(t, err)
			if got := h.Generate(message); !reflect.DeepEqual(got, messageMAC) {
				t.Errorf("MAC.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
