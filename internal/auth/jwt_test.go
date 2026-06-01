package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)	

func TestJWT(t *testing.T) {
	correctSecret, incorrectSecret := "asad-yrtv-regt-felg-jvwn", "thisIsTheWrongSecretToken"
	userID := uuid.New()

	type testUnit struct {
		Name string
		Uuid uuid.UUID
		TokenSecret string
		ExpireLimit time.Duration
	}

	correctCases := []testUnit{
		{
		Name: "General Test",
		Uuid: userID,
		TokenSecret: correctSecret,
		ExpireLimit: 3 * time.Second,
		},
		{
		Name: "Typed Out Secret Token Test",
		Uuid: userID,
		TokenSecret: "asad-yrtv-regt-felg-jvwn",
		ExpireLimit: 5 * time.Second,
		},
	}

	for _, tt := range correctCases {
		t.Run(tt.Name, func(t *testing.T) {
			jwt, err := MakeJWT(tt.Uuid, tt.TokenSecret, tt.ExpireLimit)
			if err != nil {
				t.Fatalf("error creating jwt: %q", err)
			}

			_, err = ValidateJWT(jwt, correctSecret)
			if err != nil {
				t.Fatalf("error validating jwt: %q", err)
			}
		})
	}

	incorrectCases := []testUnit {
		{
		Name: "Incorrect Secret Test",
		Uuid: userID,
		TokenSecret: incorrectSecret,
		ExpireLimit: 10 * time.Second,
		},
		{
		Name: "Expired Test",
		Uuid: userID,
		TokenSecret: correctSecret,
		ExpireLimit: -1 * time.Second,
		},
		{
		Name: "Empty Secret",
		Uuid: userID,
		TokenSecret: "",
		ExpireLimit: 3 * time.Second,
		},
		
	}

	for _, tt := range incorrectCases {
		t.Run(tt.Name, func(t *testing.T) {
			jwt, err := MakeJWT(tt.Uuid, tt.TokenSecret, tt.ExpireLimit)
			if err != nil {
				t.Fatal("unable to create jwt")
			}

			_, err = ValidateJWT(jwt, correctSecret)
			if err == nil {
				t.Fatal("incorrect test case should not pass test")
			}
		})
	}
}


