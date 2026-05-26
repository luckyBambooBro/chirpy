package auth

import (
	"testing"
)

func TestPassword(t *testing.T) {
	testCases := []struct{
		name string
		userInput string
	}{
		{
			name: "alllower",
			userInput: 	"somepassword",
		},
		{
			name: "caps",
			userInput: "SomePassword",
		},
		{
			name: "symbols",
			userInput: "!@#$%^&*()_+",
		},
		{
			name: "blank",
			userInput: "",
		},
	}
	
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			hash, err := HashPassword(testCase.userInput)
			if err != nil {
				t.Fatal(err)
			}

			match, err := CheckPasswordHash(testCase.userInput, hash)
			if err != nil {
				t.Fatal(err)
			}
			if !match {
				t.Fatalf("test case %v failed", testCase.name)
			}

			wrongMatch, _ := CheckPasswordHash("this_is_a_wrong_password", hash)
            if wrongMatch { //this should be false, therefore if true it is an error
                t.Fatalf("test case %v failed: an incorrect password successfully authenticated against the hash!", testCase.name)
            }
		})
	}
}