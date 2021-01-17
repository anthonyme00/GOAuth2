package GOAuth2

import (
	"testing"
	"time"
)

func (a OAuthAPIConnectionData) equal(b OAuthAPIConnectionData) bool {
	if a.ClientId == b.ClientId && a.ClientSecret == b.ClientSecret {
		if len(a.Scopes) == len(b.Scopes) {
			for k := range a.Scopes {
				if a.Scopes[k] != b.Scopes[k] {
					return false
				}
			}
		}
	}
	return true
}

func (a OAuth2Token) equal(b OAuth2Token) bool {
	if a.AccessToken == b.AccessToken && a.ConnectionData.equal(b.ConnectionData) && a.Expiration == b.Expiration &&
		a.IdToken == b.IdToken && a.LastRefresh == b.LastRefresh && a.RefreshToken == b.RefreshToken && a.Scope == b.Scope &&
		a.TokenType == b.TokenType {
		return true
	} else {
		return false
	}
}

func TestTokenExpiry(t *testing.T) {
	inputToken := OAuth2Token{
		Expiration:  5,
		LastRefresh: time.Now(),
	}

	if inputToken.IsExpired(0) {
		t.Errorf("Failed ! token shouldn't be expired")
		t.Fail()
	}

	timeDelta, _ := time.ParseDuration("-6s")
	inputToken.LastRefresh = time.Now().Add(timeDelta)

	if !inputToken.IsExpired(0) {
		t.Errorf("Failed ! token should be expired")
		t.Fail()
	}
}

func TestTokenSerialization(t *testing.T) {
	inputToken := OAuth2Token{
		AccessToken:    "test",
		RefreshToken:   "test",
		Expiration:     3,
		IdToken:        "test",
		Scope:          "test",
		TokenType:      "test",
		LastRefresh:    time.Now(),
		ConnectionData: OAuthAPIConnectionData{},
	}
	expectedOutput := OAuth2Token{
		AccessToken:    "test",
		RefreshToken:   "test",
		Expiration:     3,
		IdToken:        "test",
		Scope:          "test",
		TokenType:      "test",
		LastRefresh:    time.Now(),
		ConnectionData: OAuthAPIConnectionData{},
	}

	serializedToken, err := inputToken.Serialize()
	if err != nil {
		t.Errorf("Failed ! Unable to serialize token %s", err)
		t.FailNow()
	}

	newToken := OAuth2Token{}
	err = newToken.Deserialize(serializedToken)
	if err != nil {
		t.Errorf("Failed ! Unable to deserialize token %s", err)
		t.FailNow()
	}

	if newToken.equal(expectedOutput) {
		t.Log("Success !")
	}
}

func TestTokenEncryptionSerialization(t *testing.T) {
	key := []byte("SECRETKEY")

	inputToken := OAuth2Token{
		AccessToken:    "test",
		RefreshToken:   "test",
		Expiration:     3,
		IdToken:        "test",
		Scope:          "test",
		TokenType:      "test",
		LastRefresh:    time.Now(),
		ConnectionData: OAuthAPIConnectionData{},
	}
	expectedOutput := OAuth2Token{
		AccessToken:    "test",
		RefreshToken:   "test",
		Expiration:     3,
		IdToken:        "test",
		Scope:          "test",
		TokenType:      "test",
		LastRefresh:    time.Now(),
		ConnectionData: OAuthAPIConnectionData{},
	}

	serializedToken, err := inputToken.SerializeEncrypted(key)
	if err != nil {
		t.Errorf("Failed ! Unable to serialize token %s", err)
		t.FailNow()
	}

	newToken := OAuth2Token{}
	err = newToken.DeserializeEncrypted(serializedToken, key)
	if err != nil {
		t.Errorf("Failed ! Unable to deserialize token %s", err)
		t.FailNow()
	}

	if newToken.equal(expectedOutput) {
		t.Log("Success !")
	}
}

type slice []interface{}

func (a slice) compareSlice(b []interface{}) bool {
	if len(a) != len(b) {
		return false
	} else {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}

func TestScope(t *testing.T) {
	input := OAuth2Token{
		Scope: "a b c d",
	}

	expectedOutput := slice{"a", "b", "c", "d"}

	inputSlice := make([]interface{}, len(input.Scopes()))
	for i, v := range input.Scopes() {
		inputSlice[i] = v
	}

	if slice(inputSlice).compareSlice(expectedOutput) {
		t.Log("Success !")
	} else {
		t.Errorf("Failed to get scopes from token")
		t.FailNow()
	}
}
