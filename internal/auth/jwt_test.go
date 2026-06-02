package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)	

func TestJWT(t *testing.T) {
	defaultSecret := "asad-yrtv-regt-felg-jvwn"
	userID := uuid.New()

	tests := []struct {
		Name string
		Uuid uuid.UUID
		TokenSecret string
		ExpireLimit time.Duration
		WantErr bool
	}{
		{
		Name: "General Test",
		Uuid: userID,
		TokenSecret: defaultSecret,
		ExpireLimit: 3 * time.Second,
		},
		{
		Name: "Typed Out Secret Token Test",
		Uuid: userID,
		TokenSecret: "asad-yrtv-regt-felg-jvwn",
		ExpireLimit: 5 * time.Second,
		},
		{
		Name: "Incorrect Secret Test",
		Uuid: userID,
		TokenSecret: "incorrectSecret",
		ExpireLimit: 10 * time.Second,
		WantErr: true,
		},
		{
		Name: "Expired Test",
		Uuid: userID,
		TokenSecret: defaultSecret,
		ExpireLimit: -1 * time.Second,
		WantErr: true,
		},
		{
		Name: "Empty Secret",
		Uuid: userID,
		TokenSecret: "",
		ExpireLimit: 3 * time.Second,
		WantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			jwt, err := MakeJWT(tt.Uuid, tt.TokenSecret, tt.ExpireLimit)
			if err != nil {
				t.Fatalf("unable to create jwt: %v", err)
			}

			_, err = ValidateJWT(jwt, defaultSecret)
			if (err != nil) != tt.WantErr {
				t.Fatalf("validateJWT err: %v\nWantErr: %v", err, tt.WantErr)
			}
		})
	}
}


