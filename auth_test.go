package myanimelist

import (
	"os"
	"testing"
)

func TestMAL_RefreshToken(t *testing.T) {
	if _, ok := os.LookupEnv("TRAVIS"); ok {
		t.Skip("we want to run this test only locally")
	}
	mal := ExampleMAL

	if mal.Auth.clientID == "" {
		t.Fatal("you need to set clientID in your secret.yaml for this test")
	}
	if mal.Auth.clientSecret == "" {
		t.Fatal("you need to set clientSecret in your secret.yaml for this test")
	}
	if mal.Auth.refreshToken == "" {
		t.Fatal("you need to set refreshToken in your secret.yaml for this test")
	}

	userInfo, err := mal.Auth.RefreshToken()
	if err != nil {
		t.Fatalf("refreshToken got error: %v", err)
	}

	if userInfo.AccessToken == "" || userInfo.RefreshToken == "" {
		t.Fatal("Got empty fields")
	}

	// if everything went good, we need to update our credentials in file for future usage
	secretFileWrite(&TestCredentials{
		AccessToken:  userInfo.AccessToken,
		RefreshToken: userInfo.RefreshToken,
	})

	// and in example structure if we keep doing tests
	ExampleMAL.Auth.userToken = userInfo.AccessToken
	ExampleMAL.Auth.refreshToken = userInfo.RefreshToken
}
