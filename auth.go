package myanimelist

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Auth contain all authorization-related data
type Auth struct {
	mal *MAL

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

var authorizeEndpoint = "https://myanimelist.net/v1/oauth2/authorize"
var tokenEndpoint = "https://myanimelist.net/v1/oauth2/token"

// LoginURL starts OAuth process and return login URL.
// For additional info use this: https://myanimelist.net/apiconfig/references/authorization.
func (a *Auth) LoginURL() string {
	// Generate PKCE codes - https://tools.ietf.org/html/rfc7636
	a.codeVerifier = codeVerifier()
	a.codeChallenge = codeChallenge(a.codeVerifier, codeChallengePlain)

	reqURL, _ := url.Parse(authorizeEndpoint)

	q := reqURL.Query()
	q.Set("response_type", "code")
	q.Set("client_id", a.clientID)
	q.Set("redirect_uri", a.redirectURL)
	q.Set("code_challenge", a.codeChallenge)

	reqURL.RawQuery = q.Encode()

	return reqURL.String()
}

// RetrieveToken use received from user's authorization code and send
// it to server to receive user access token
func (a *Auth) ExchangeToken(authCode string) (*UserCredentials, error) {
	method := http.MethodPost
	path := tokenEndpoint
	data := url.Values{
		"client_id":     {a.clientID},
		"client_secret": {a.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {authCode},
		"redirect_uri":  {a.redirectURL},
		"code_verifier": {a.codeVerifier},
	}

	tokenResp := new(tokenResponse)
	if err := a.mal.request(tokenResp, method, path, data); err != nil {
		return nil, err
	}

	expirationDuration, err := time.ParseDuration(fmt.Sprintf("%ds", tokenResp.ExpiresIn))
	if err != nil {
		a.mal.logger.Printf("Parse duration error: %s\n", err)
		return nil, err
	}

	expireAt := time.Now().Add(expirationDuration)

	user := &UserCredentials{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpireAt:     expireAt,
	}
	a.SetTokenInfo(user.AccessToken, user.RefreshToken, user.ExpireAt)

	return user, nil
}

type codeChallengeMethod string

const (
	codeChallengePlain codeChallengeMethod = "plain"
	// redundant for now, because MyAnimeList support only plain method yet
	//codeChallengeSHA256 codeChallengeMethod = "S256"
)

// codeVerifier generates random string, as RFC7636 require.
// Reference: https://tools.ietf.org/html/rfc7636#section-4.1.
func codeVerifier() string {
	var allowedSymbols = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~")

	// Minimum - 43, maximum is 128 symbols
	return randomString(128, allowedSymbols)
}

// codeChallenge encode our generated string. For additional info look at
// Section 4.2 of RFC7636 - https://tools.ietf.org/html/rfc7636#section-4.2.
func codeChallenge(code string, _ codeChallengeMethod) string {
	// redundant because MyAnimeList support only plain method yet
	//if mode == codeChallengePlain {
	return code
	//}
	// BASE64URL-ENCODE(SHA256(ASCII(code_verifier)))
	//s := sha256.New()
	//s.Write([]byte(code))
	//encoded := s.Sum(nil)
	//
	//return base64.URLEncoding.EncodeToString(encoded)
}

func (a *Auth) RefreshToken() (*UserCredentials, error) {
	method := http.MethodPost
	path := tokenEndpoint
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {a.refreshToken},
		"client_id":     {a.clientID},
		"client_secret": {a.clientSecret},
	}

	tokenResp := new(tokenResponse)
	if err := a.mal.request(tokenResp, method, path, data); err != nil {
		return nil, err
	}

	expirationDuration, err := time.ParseDuration(fmt.Sprintf("%ds", tokenResp.ExpiresIn))
	if err != nil {
		a.mal.logger.Printf("Parse duration error: %s\n", err)
		return nil, err
	}

	expireAt := time.Now().Add(expirationDuration)

	user := &UserCredentials{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpireAt:     expireAt,
	}

	a.SetTokenInfo(user.AccessToken, user.RefreshToken, user.ExpireAt)
	return user, nil
}

type UserCredentials struct {
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
}

// GetTokenInfo returns all required user's credentials: access token,
// refresh token and access token's expiration date.
func (a *Auth) GetTokenInfo() *UserCredentials {
	return &UserCredentials{
		AccessToken:  a.userToken,
		RefreshToken: a.refreshToken,
		ExpireAt:     a.tokenExpireAt,
	}
}

// SetTokenInfo completely rewrites saved user's credentials, so use it very careful.
// In case you erased correct tokens - lead user to authorization page again.
func (a *Auth) SetTokenInfo(accessToken string, refreshToken string, expire time.Time) {
	a.userToken = accessToken
	a.refreshToken = refreshToken
	a.tokenExpireAt = expire
}

type tokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
