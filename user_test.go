package myanimelist

import (
	"testing"
)

func TestMAL_UserInformation(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Example",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.UserInformation()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("UserInformation() Got nil")
				return
			}
			if got.ID == 0 {
				t.Error("UserInformation() Got user with empty ID")
				return
			}
		})
	}
}
