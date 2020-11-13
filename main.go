// myanimelist is a small library to simplify usege of MyAnimeList's API
package myanimelist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var apiEndpoint = "https://api.myanimelist.net/v2/"

type MAL struct {
	// Host contain API entry point
	host string

	client *http.Client

	logger *log.Logger

	// Auth contain all authorization-related data
	auth auth
}

// Auth contain all authorization-related data
type auth struct {
	// application credentials required for authorization
	clientID, clientSecret string

	// token to identify user. Required for every request
	userToken string

	// token expiration time
	tokenExpireAt time.Time

	// required for receiving new user token
	refreshToken string

	// part of RFC7636 authorization
	codeVerifier, codeChallenge string

	// url to redirect after myAnimeList authorization
	redirectURL string
}

// New creates new MyAnimeList client with specified parameters.
// Every api method require authorization so you need to provide
// all auth-related data before client's initialisation.
// If you plan to set user's tokens manually with SetTokenInfo(),
// instead of authorization - you can specify "/" as redirect URL.
func New(config Config) (*MAL, error) {
	if config.ClientID == "" {
		return nil, errors.New("field ClientID is required")
	}
	if config.ClientSecret == "" {
		return nil, errors.New("field ClientSecret is required")
	}
	if config.RedirectURL == "" {
		return nil, errors.New("field RedirectURL is required")
	}

	mal := &MAL{
		host: apiEndpoint,
		auth: auth{
			clientID:     config.ClientID,
			clientSecret: config.ClientSecret,
			redirectURL:  config.RedirectURL,
		},
		client: &http.Client{Timeout: 5 * time.Second},
		logger: log.New(os.Stderr, "[MAL] ", 0),
	}

	if config.HTTPClient != nil {
		mal.client = config.HTTPClient
	}

	if config.Logger != nil {
		mal.logger = config.Logger
	}

	return mal, nil
}

// Config stores important data to create new MyAnimeList client.
// HTTPClient and Logger is optional.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	HTTPClient   *http.Client
	Logger       *log.Logger
}

type ErrorResponse struct {
	Err     string `json:"error"`
	Message string `json:"message, omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("myanimelist returned error: %s. With message: %s", e.Err, e.Message)
}

// GetTokenInfo returns all required user's credentials: access token,
// refresh token and access token's expiration date.
func (mal *MAL) GetTokenInfo() *UserCredentials {
	return &UserCredentials{
		AccessToken:  mal.auth.userToken,
		RefreshToken: mal.auth.refreshToken,
		ExpireAt:     mal.auth.tokenExpireAt,
	}
}

// SetTokenInfo completely rewrites saved user's credentials, so use it very careful.
// In case you erased correct tokens - lead user to authorization page again.
func (mal *MAL) SetTokenInfo(accessToken string, refreshToken string, expire time.Time) {
	mal.auth.userToken = accessToken
	mal.auth.refreshToken = refreshToken
	mal.auth.tokenExpireAt = expire
}

// requestRaw makes actual request and returns everything we got
func (mal *MAL) requestRaw(method string, path string, data url.Values) (*http.Response, error) {
	var body = new(strings.Reader)

	baseURL, _ := url.Parse(mal.host)

	apiURL, err := baseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	log.Printf("path: %s", apiURL.String())

	if len(data) > 0 {
		switch method {
		case http.MethodGet:
			// append query
			apiURL.RawQuery += "&" + data.Encode()
		default:
			body = strings.NewReader(data.Encode())
		}
	}

	req, err := http.NewRequest(method, apiURL.String(), body)
	if err != nil {
		return nil, err
	}

	// Only add authorization header if its not auth-related request
	if !strings.Contains(apiURL.Path, "v1/oauth2") {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mal.auth.userToken))
	}

	if method == http.MethodPost || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(int(body.Size())))
	}

	resp, err := mal.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// request is small wrapper around requestRaw to avoid multiple ReadBody->Unmarshal->CloseBody chains
func (mal *MAL) request(destination interface{}, method string, path string, data url.Values) error {
	resp, err := mal.requestRaw(method, path, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		if err := json.Unmarshal(respBody, destination); err != nil {
			return err
		}
		return nil
	}

	// Try to parse error message
	errorMsg := new(ErrorResponse)
	if err := json.Unmarshal(respBody, errorMsg); err != nil {
		return err
	}
	return errorMsg
}

type Paging struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}

// PagingSettings contains limit and offset fields, which are applicable for almost every request.
// Also, im general, max && default Limit value is 100.
// But for actual information refer to certain method's official documentation.
type PagingSettings struct {
	Limit  int
	Offset int
}

func (s *PagingSettings) set(values *url.Values) {
	if s.Limit != 0 {
		values.Set("limit", strconv.Itoa(s.Limit))
	}
	if s.Offset != 0 {
		values.Set("offset", strconv.Itoa(s.Offset))
	}
}

func (mal *MAL) getPage(result interface{}, p Paging, direction int8, limit []int) error {
	var pageURL string

	if direction < 0 {
		pageURL = p.Previous
	} else {
		pageURL = p.Next
	}

	if pageURL == "" {
		return errors.New("there is no more pages")
	}

	if len(limit) > 0 {
		if limit[0] > 0 {
			pageObj, err := url.Parse(pageURL)
			if err != nil {
				return fmt.Errorf("something wrong with url: %s", err)
			}

			newQuery := pageObj.Query()
			newQuery.Set("limit", fmt.Sprint(limit[0]))
			pageObj.RawQuery = newQuery.Encode()

			pageURL = pageObj.String()
		}
	}

	return mal.request(result, http.MethodGet, pageURL, url.Values{})
}