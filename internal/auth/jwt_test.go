package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)	

func TestJWT(t *testing.T) {
	correctSecret, incorrectSecret := "asad-yrtv-regt-felg-jvwn", "thisIsTheWrongSecretToken"

	type testUnit struct {
		Name string
		Uuid uuid.UUID
		TokenSecret string
		ExpireLimit time.Duration
	}

	correctCases := []testUnit{
		{
		Name: "General Test",
		Uuid: uuid.New(),
		TokenSecret: correctSecret,
		ExpireLimit: 3 * time.Second,
		},
		{
		Name: "Typed Out Secret Token Test",
		Uuid: uuid.New(),
		TokenSecret: "asad-yrtv-regt-felg-jvwn",
		ExpireLimit: 5 * time.Second,
		},
	}

	for _, testCase := range correctCases {
		t.Run(testCase.Name, func(t *testing.T) {
			jwt, err := MakeJWT(testCase.Uuid, testCase.TokenSecret, testCase.ExpireLimit)
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
		Uuid: uuid.New(),
		TokenSecret: incorrectSecret,
		ExpireLimit: 10 * time.Second,
		},
		{
		Name: "Expired Test",
		Uuid: uuid.New(),
		TokenSecret: correctSecret,
		ExpireLimit: -1 * time.Second,
		},
		{
		Name: "Empty Secret",
		Uuid: uuid.New(),
		TokenSecret: "",
		ExpireLimit: 3 * time.Second,
		},
		
	}

	for _, testCase := range incorrectCases {
		t.Run(testCase.Name, func(t *testing.T) {
			jwt, err := MakeJWT(testCase.Uuid, testCase.TokenSecret, testCase.ExpireLimit)
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


