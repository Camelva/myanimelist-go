package myanimelist

import (
	"testing"
)

func TestMAL_User_Info(t *testing.T) {
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
			got, err := ExampleMAL.User.Info()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_User_Info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestMAL_User_Info() Got nil")
				return
			}
			if got.ID == 0 {
				t.Error("TestMAL_User_Info() Got user with empty ID")
				return
			}
		})
	}
}
