package myanimelist

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type tokenResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var authorizeEndpoint = "https://myanimelist.net/v1/oauth2/authorize"
var tokenEndpoint = "https://myanimelist.net/v1/oauth2/token"

// AuthURL starts OAuth process and return user's auth URL.
// For additional info use this: https://myanimelist.net/apiconfig/references/authorization.
func (mal *MAL) AuthURL() string {
	// Generate PKCE codes - https://tools.ietf.org/html/rfc7636
	mal.auth.codeVerifier = codeVerifier()
	mal.auth.codeChallenge = codeChallenge(mal.auth.codeVerifier, codeChallengePlain)

	reqURL, _ := url.Parse(authorizeEndpoint)

	q := reqURL.Query()
	q.Set("response_type", "code")
	q.Set("client_id", mal.auth.clientID)
	q.Set("redirect_uri", mal.auth.redirectURL)
	q.Set("code_challenge", mal.auth.codeChallenge)

	reqURL.RawQuery = q.Encode()

	return reqURL.String()
}

// RetrieveToken use received from user's authorization code and send
// it to server to receive user access token
func (mal *MAL) ExchangeToken(authCode string) (*UserCredentials, error) {
	method := http.MethodPost
	path := tokenEndpoint
	data := url.Values{
		"client_id":     {mal.auth.clientID},
		"client_secret": {mal.auth.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {authCode},
		"redirect_uri":  {mal.auth.redirectURL},
		"code_verifier": {mal.auth.codeVerifier},
	}

	tokenResp := new(tokenResponse)
	if err := mal.request(tokenResp, method, path, data); err != nil {
		return nil, err
	}

	expirationDuration, err := time.ParseDuration(fmt.Sprintf("%ds", tokenResp.ExpiresIn))
	if err != nil {
		mal.logger.Printf("Parse duration error: %s\n", err)
		return nil, err
	}

	expireAt := time.Now().Add(expirationDuration)

	user := &UserCredentials{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpireAt:     expireAt,
	}
	mal.SetTokenInfo(user.AccessToken, user.RefreshToken, user.ExpireAt)

	return user, nil
}

type codeChallengeMethod string

const (
	codeChallengePlain  codeChallengeMethod = "plain"
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

func (mal *MAL) RefreshToken() (*UserCredentials, error) {
	method := http.MethodPost
	path := tokenEndpoint
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {mal.auth.refreshToken},
		"client_id":     {mal.auth.clientID},
		"client_secret": {mal.auth.clientSecret},
	}

	tokenResp := new(tokenResponse)
	if err := mal.request(tokenResp, method, path, data); err != nil {
		return nil, err
	}

	expirationDuration, err := time.ParseDuration(fmt.Sprintf("%ds", tokenResp.ExpiresIn))
	if err != nil {
		mal.logger.Printf("Parse duration error: %s\n", err)
		return nil, err
	}

	expireAt := time.Now().Add(expirationDuration)

	user := &UserCredentials{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpireAt:     expireAt,
	}

	mal.SetTokenInfo(user.AccessToken, user.RefreshToken, user.ExpireAt)
	return user, nil
}

type UserCredentials struct {
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
}
