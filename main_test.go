package myanimelist

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var ExampleMAL *MAL

var secretFileName = "secret.yaml"

type TestCredentials struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	AccessToken  string `yaml:"accessToken"`
	RefreshToken string `yaml:"refreshToken"`
}

func init() {
	secretData := secretFileRead()
	if secretData.AccessToken == "" {
		log.Fatalln("access token is required to run any test")
	}

	testClient, err := New(Config{
		ClientID:     secretData.ClientID,
		ClientSecret: secretData.ClientSecret,
		RedirectURL:  "/",
		HTTPClient:   &http.Client{Timeout: 5 * time.Second},
		Logger:       log.New(os.Stderr, "[TEST MAL]", 0),
	})
	if err != nil {
		log.Fatalf("can't init testClient: %s", err)
	}

	testClient.auth.userToken = secretData.AccessToken
	testClient.auth.refreshToken = secretData.RefreshToken

	ExampleMAL = testClient
}

func secretFileRead() *TestCredentials {
	secretFilePath := filepath.Join("testdata", secretFileName)
	secretFileContent, err := ioutil.ReadFile(secretFilePath)
	if err != nil {
		log.Fatalln("make sure to create secret.yaml file inside testdata folder before running tests")
	}

	secretData := new(TestCredentials)
	if err := yaml.Unmarshal(secretFileContent, secretData); err != nil {
		log.Fatalf("can't parse config file: %v\n", err)
	}
	return secretData
}

// secretFileWrite updates content of your file with credentials.
// Changes only provided values, keeping others original.
func secretFileWrite(config *TestCredentials) {
	secretData := secretFileRead()
	if config.ClientSecret != "" {
		secretData.ClientSecret = config.ClientSecret
	}

	if config.ClientID != "" {
		secretData.ClientID = config.ClientID
	}

	if config.AccessToken != "" {
		secretData.AccessToken = config.AccessToken
	}

	if config.RefreshToken != "" {
		secretData.RefreshToken = config.RefreshToken
	}
	newSecretData, err := yaml.Marshal(secretData)
	if err != nil {
		log.Fatalf("can't update secret file, error while marshaling new data: %s", err)
	}

	secretFilePath := filepath.Join("testdata", secretFileName)
	if err := ioutil.WriteFile(secretFilePath, newSecretData, 0644); err != nil {
		log.Fatalf("can't rewrite secret file: %s", err)
	}
}
