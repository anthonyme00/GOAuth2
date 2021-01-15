package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	httplistener "github.com/anthonyme00/GOAuth2/httplistener"
	"github.com/anthonyme00/GOAuth2/util"
)

// https://developers.google.com/identity/protocols/oauth2/scopes
type OAuthAPIConnectionData struct {
	ClientId     string
	ClientSecret string
	Scopes       []string
}

type OAuth2Token struct {
	AccessToken    string  `json:"access_token"`
	RefreshToken   string  `json:"refresh_token"`
	Expiration     float64 `json:"expires_in"`
	IdToken        string  `json:"id_token"`
	Scope          string  `json:"scope"`
	TokenType      string  `json:"token_type"`
	LastRefresh    time.Time
	ConnectionData OAuthAPIConnectionData
}

func (token *OAuth2Token) IsExpired(expirationThreshold float64) bool {
	currentTime := time.Now()

	timeDelta, _ := time.ParseDuration(fmt.Sprintf("%vs", token.Expiration-expirationThreshold))
	if token.LastRefresh.Add(timeDelta).Before(currentTime) {
		return true
	} else {
		return false
	}
}

func (token *OAuth2Token) Scopes() []string {
	return strings.Split(token.Scope, " ")
}

func (token *OAuth2Token) GetAccessToken(expirationThreshold float64) string {
	if token.IsExpired(expirationThreshold) == false {
		return token.AccessToken
	} else {
		token.Refresh()
		return token.AccessToken
	}
}

func (token *OAuth2Token) Refresh() {
	refreshQuery := util.UrlQuery{
		"client_id":     token.ConnectionData.ClientId,
		"client_secret": token.ConnectionData.ClientSecret,
		"grant_type":    "refresh_token",
		"refresh_token": token.RefreshToken,
	}

	resp, err := http.Post(OAuth2TokenEndpoint+refreshQuery.CreateQuery(), "application/x-www-form-urlencoded", nil)

	if err == nil {
		responseBody, _ := ioutil.ReadAll(resp.Body)

		resp.Body.Close()

		json.Unmarshal(responseBody, token)
		token.LastRefresh = time.Now()
	}
}

func (token *OAuth2Token) Serialize() ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(token)

	if err != nil {
		return nil, err
	} else {
		return buffer.Bytes(), nil
	}
}

func (token *OAuth2Token) SerializeEncrypted(key []byte) ([]byte, error) {
	serializedToken, err := token.Serialize()
	if err != nil {
		return nil, err
	}

	encryptedToken, err := util.Encrypt(serializedToken, key)
	if err != nil {
		return nil, err
	}

	return encryptedToken, nil
}

func (token *OAuth2Token) Deserialize(data []byte) error {
	dataBuffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(dataBuffer)

	err := decoder.Decode(token)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (token *OAuth2Token) DeserializeEncrypted(data []byte, key []byte) error {
	decryptedToken, err := util.Decrypt(data, key)
	if err != nil {
		return err
	}

	token.Deserialize(decryptedToken)

	return nil
}

var OAuth2AuthEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
var OAuth2TokenEndpoint = "https://oauth2.googleapis.com/token"

func GetOAuth2Token(data OAuthAPIConnectionData) (*OAuth2Token, error) {
	code_verifier := util.GenerateBase64URLnopadding(64)
	code_challenge := util.GenerateSHA256(code_verifier)

	authorizationListener := httplistener.HTTPListener{}
	authorizationListenerHandle := authorizationListener.OpenListener(httplistener.RedirectToURL("https://google.com/"))
	authorizationRedirect := authorizationListener.GetUrl()

	//https://developers.google.com/identity/protocols/oauth2/native-app#step-2:-send-a-request-to-googles-oauth-2.0-server
	requestAuthorizationQuery := util.UrlQuery{
		"client_id":             data.ClientId,
		"response_type":         "code",
		"redirect_uri":          url.QueryEscape(authorizationRedirect),
		"scope":                 url.QueryEscape(strings.Join(data.Scopes, " ")),
		"code_challenge":        util.Base64URLNoPadding(code_challenge[:]),
		"code_challenge_method": "S256",
	}

	httplistener.OpenInBrowser(OAuth2AuthEndpoint + requestAuthorizationQuery.CreateQuery())

	authorizationResponse := authorizationListener.GetResponse(authorizationListenerHandle)
	authorizationResponseQueries := authorizationResponse.URL.Query()

	if val, ok := authorizationResponseQueries["code"]; ok {
		tokenRequestQuery := util.UrlQuery{
			"redirect_uri":  url.QueryEscape(authorizationRedirect),
			"client_id":     data.ClientId,
			"client_secret": data.ClientSecret,
			"code":          val[0],
			"code_verifier": code_verifier,
			"grant_type":    "authorization_code",
		}

		resp, err := http.Post(OAuth2TokenEndpoint+tokenRequestQuery.CreateQuery(), "application/x-www-form-urlencoded", nil)

		if err == nil {
			authorizationToken := OAuth2Token{}
			responseBody, _ := ioutil.ReadAll(resp.Body)

			resp.Body.Close()

			json.Unmarshal(responseBody, &authorizationToken)
			authorizationToken.LastRefresh = time.Now()
			authorizationToken.ConnectionData = data

			return &authorizationToken, nil
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("Error getting authorization: " + authorizationResponseQueries["error"][0])
	}
}
