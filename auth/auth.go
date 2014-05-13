package auth

import (
	"code.google.com/p/goauth2/oauth"
	"errors"
	"fmt"
	"github.com/prasmussen/gdrive/util"
	"net/http"
)

// Get auth code from user
func promptUserForAuthCode(config *oauth.Config) string {
	authUrl := config.AuthCodeURL("state")
	fmt.Println("Go to the following link in your browser:")
	fmt.Printf("%v\n\n", authUrl)
	return util.Prompt("Enter verification code: ")
}

// Returns true if we have a valid cached token
func hasValidToken(cacheFile oauth.CacheFile, transport *oauth.Transport) bool {
	// Check if we have a cached token
	token, err := cacheFile.Token()
	if err != nil {
		return false
	}

	// Refresh token if its expired
	if token.Expired() {
		transport.Token = token
		err = transport.Refresh()
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return true
}

func GetOauth2Client(clientId, clientSecret, cachePath string, promptUser bool) (*http.Client, error) {
	cacheFile := oauth.CacheFile(cachePath)

	config := &oauth.Config{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Scope:        "https://www.googleapis.com/auth/drive",
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		TokenCache:   cacheFile,
	}

	transport := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
	}

	// Return client if we have a valid token
	if hasValidToken(cacheFile, transport) {
		return transport.Client(), nil
	}

	if !promptUser {
		return nil, errors.New("no valid token found")
	}

	// Get auth code from user and request a new token
	code := promptUserForAuthCode(config)
	_, err := transport.Exchange(code)
	if err != nil {
		return nil, err
	}
	return transport.Client(), nil
}
