# GOAuth2

[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/anthonyme00/GOAuth2.svg)](https://pkg.go.dev/github.com/anthonyme00/GOAuth2)
[![Coverage Status](https://coveralls.io/repos/github/anthonyme00/GOAuth2/badge.svg?branch=master)](https://coveralls.io/github/anthonyme00/GOAuth2?branch=master)

A simple gog library for getting Google OAuth2 token 

example usage:

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
}
```