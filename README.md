# GOAuth2

[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/anthonyme00/GOAuth2.svg)](https://pkg.go.dev/github.com/anthonyme00/GOAuth2)
[![Coverage Status](https://coveralls.io/repos/github/anthonyme00/GOAuth2/badge.svg?branch=master)](https://coveralls.io/github/anthonyme00/GOAuth2?branch=master)
![Travis (.com)](https://img.shields.io/travis/com/anthonyme00/GOAuth2)

A simple go library for getting Google OAuth2 token 

### Examples
***
Using :

```
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	oauth "github.com/anthonyme00/GOAuth2"
)

func main() {
	connectionData := oauth.OAuthAPIConnectionData{
		ClientId:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		Scopes:       []string{"https://www.googleapis.com/auth/drive.file"},
	}

	token, err := oauth.GetOAuth2Token(connectionData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", token)

	//Using token in a query
	query := UrlQuery { 
		"access_token":	token.GetAccessToken(),
		"other_query":	"aquery",
	}
	DummyGoogleAPI.Call(query)
}
```

Refreshing tokens :

```
token, err := oauth.GetOAuth2Token(connectionData)

//Automatic token refreshing when getting access token
//3 seconds threshold before actual expiration of token
current_token := token.GetAccessToken(3))

//Manual token refreshing
token.Refresh()

//Checking if token is expired
if token.IsExpired(0) {
	fmt.Println("My token is expired...")
}
```

Serializing and saving to a file :

```
key := []byte("SUPER_SECRET_KEY")

//Save Token
//You can use the non encrypted version with token.Serialize()
encryptedToken, _ := token.SerializeEncrypted(key)
ioutil.WriteFile("mytoken.dat", encryptedToken, os.FileMode(os.O_CREATE|os.O_TRUNC))

//Load Token
//To load non encrypted token, use token.Deserialize()
tokenToDecrypt, _ := ioutil.ReadFile("mytoken.dat")
decryptedToken := oauth.OAuth2Token{}
decryptedToken.DeserializeEncrypted(tokenToDecrypt, key)
fmt.Printf("%+v\n", token)
```

### Misc
---

Google API Scopes : https://developers.google.com/identity/protocols/oauth2/scopes

Google OAuth2 Docs : https://developers.google.com/identity/protocols/oauth2