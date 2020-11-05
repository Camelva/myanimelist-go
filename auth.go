package myanimelist

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/dchest/uniuri"
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
	mal.auth.codeChallenge = codeChallenge(mal.auth.codeVerifier, CodeChallengePlain)

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
	CodeChallengePlain  codeChallengeMethod = "plain"
	CodeChallengeSHA256 codeChallengeMethod = "S256"
)

// codeVerifier generates random string, as RFC7636 require.
// Reference: https://tools.ietf.org/html/rfc7636#section-4.1.
func codeVerifier() string {
	var allowedSymbols = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~")

	// Minimum - 43, maximum is 128 symbols
	return uniuri.NewLenChars(128, allowedSymbols)
}

// codeChallenge encode our generated string. For additional info look at
// Section 4.2 of RFC7636 - https://tools.ietf.org/html/rfc7636#section-4.2.
func codeChallenge(code string, mode codeChallengeMethod) string {
	if mode != CodeChallengeSHA256 {
		return code
	}

	// BASE64URL-ENCODE(SHA256(ASCII(code_verifier)))
	s := sha256.New()
	s.Write([]byte(code))
	encoded := s.Sum(nil)

	return base64.URLEncoding.EncodeToString(encoded)
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

//func (mal *MAL) customRequest(destination interface{}, method string, path string, data url.Values) error {
//	body := strings.NewReader(data.Encode())
//	req, err := http.NewRequest(method, path, body)
//	if err != nil {
//		return err
//	}
//	if method == http.MethodPost {
//		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//		req.Header.Add("Content-Length", strconv.Itoa(int(body.Size())))
//	}
//
//	reqDump, _ := httputil.DumpRequestOut(req, true)
//	mal.logger.Printf("====== REQUEST: =======\n%s\n", string(reqDump))
//
//	resp, err := mal.client.Do(req)
//
//	respDump, _ := httputil.DumpResponse(resp, true)
//	mal.logger.Printf("====== RESPOSNE: =======\n%s\n", string(respDump))
//
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	respBody, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return err
//	}
//
//	if resp.StatusCode == http.StatusOK {
//		if err := json.Unmarshal(respBody, destination); err != nil {
//			return err
//		}
//		return nil
//	}
//
//	// Try to parse error message
//	errorMsg := new(ErrorResponse)
//	if err := json.Unmarshal(respBody, errorMsg); err != nil {
//		return err
//	}
//	return errorMsg
//}

type UserCredentials struct {
	AccessToken  string
	RefreshToken string
	ExpireAt     time.Time
}
