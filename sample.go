package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	key := []byte("SUPER_SECRET_KEY")

	//Create Token
	oauthdata := OAuthAPIConnectionData{
		ClientId:     "CLIENT_ID_FROM_GOOGLE_DASHBOARD",
		ClientSecret: "CLIENT_SECRET_FROM_GOOGLE_DASHBOARD",
		Scopes:       []string{"https://www.googleapis.com/auth/drive.file"},
	}

	token, _ := GetOAuth2Token(oauthdata)
	fmt.Printf("%+v\n", token)

	//Save Token
	encryptedToken, _ := token.SerializeEncrypted(key)
	ioutil.WriteFile("mytoken.dat", encryptedToken, os.FileMode(os.O_CREATE|os.O_TRUNC))

	//Load Token
	tokenToDecrypt, _ := ioutil.ReadFile("mytoken.dat")
	decryptedToken := OAuth2Token{}

	decryptedToken.DeserializeEncrypted(tokenToDecrypt, key)
	fmt.Printf("%+v\n", token)
}
