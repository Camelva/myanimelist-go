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

var ExampleMAL = new(MAL)

var secretFileName = "secret.yaml"

type TestCredentials struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
	AccessToken  string `yaml:"accessToken"`
	RefreshToken string `yaml:"refreshToken"`
}

func init() {
	data := readEnv()
	secretData := secretFileRead()
	if secretData.ClientID != "" {
		data.ClientID = secretData.ClientID
	}
	if secretData.ClientSecret != "" {
		data.ClientSecret = secretData.ClientSecret
	}
	if secretData.AccessToken != "" {
		data.AccessToken = secretData.AccessToken
	}
	if secretData.RefreshToken != "" {
		data.RefreshToken = secretData.RefreshToken
	}

	if data.AccessToken == "" {
		log.Fatalln("access token is required to run any test")
	}
	if data.ClientID == "" {
		data.ClientID = "mock"
	}
	if data.ClientSecret == "" {
		data.ClientSecret = "mock"
	}

	testClient, err := New(Config{
		ClientID:     data.ClientID,
		ClientSecret: data.ClientSecret,
		RedirectURL:  "/",
		HTTPClient:   &http.Client{Timeout: 5 * time.Second},
		Logger:       log.New(os.Stderr, "[TEST MAL]", 0),
	})
	if err != nil {
		log.Fatalf("can't init testClient: %s", err)
	}

	testClient.Auth.userToken = data.AccessToken
	testClient.Auth.refreshToken = data.RefreshToken

	ExampleMAL = testClient
}

func secretFileRead() *TestCredentials {
	secretData := new(TestCredentials)

	secretFilePath := filepath.Join("testdata", secretFileName)
	secretFileContent, err := ioutil.ReadFile(secretFilePath)

	if err != nil {
		//log.Fatalln("make sure to create secret.yaml file inside testdata folder before running tests")
		return secretData
	}

	if err := yaml.Unmarshal(secretFileContent, secretData); err != nil {
		log.Printf("can't parse config file: %v\n", err)
		return secretData
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

func readEnv() *TestCredentials {
	var storage = new(TestCredentials)
	if id, ok := os.LookupEnv("CLIENT_ID"); ok {
		storage.ClientID = id
	}
	if secret, ok := os.LookupEnv("CLIENT_SECRET"); ok {
		storage.ClientSecret = secret
	}
	if accToken, ok := os.LookupEnv("ACCESS_TOKEN"); ok {
		storage.AccessToken = accToken
	}
	if refToken, ok := os.LookupEnv("REFRESH_TOKEN"); ok {
		storage.RefreshToken = refToken
	}
	return storage
}
