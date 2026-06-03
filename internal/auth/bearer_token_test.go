package auth

import (
	"net/http"
	"testing"
)

func TestAuthentication(t *testing.T) {
	cases := []struct{
		Name string
		Header http.Header
		WantBearerTokenMsg string
		WantErr bool
	}{
		{
			Name: "Valid Test",
			Header: http.Header{"Authorization": []string{"Bearer 123"}},
			WantBearerTokenMsg: "123",
			WantErr: false,
		},
		{
			Name: "Invalid Key Name Test",
			Header: http.Header{"Auth": []string{"Bearer 123"}},
			WantBearerTokenMsg: "",
			WantErr: true,
		},
		{
			Name: "Invalid String Length Test",
			Header: http.Header{"Authorization": []string{"Bearer 123 abc"}},
			WantBearerTokenMsg: "",
			WantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := GetBearerToken(tt.Header)
			
			if (err != nil) != tt.WantErr {
				t.Fatalf("expected error: %v, got: %v", tt.WantErr, err)
			}

			if tt.WantBearerTokenMsg != got {
				t.Fatalf("expected %q, got %q", tt.WantBearerTokenMsg, got)
			}
		})
	}
}