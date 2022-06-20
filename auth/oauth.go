package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authCodeFn func(*oauth2.Config, string, string) (func() (string, error), error)

func NewFileSourceClient(clientId, clientSecret, tokenFile string, authFn authCodeFn) (*http.Client, error) {
	conf := getConfig(clientId, clientSecret)

	// Read cached token
	token, exists, err := ReadToken(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to read token: %s", err)
	}

	// Require auth code if token file does not exist
	// or refresh token is missing
	if !exists || token.RefreshToken == "" {
		state, err := makeState()
		if err != nil {
			return nil, fmt.Errorf("could not build state string: %s", err)
		}
		verifier, challenge, err := makeCodeChallenge()
		if err != nil {
			return nil, fmt.Errorf("could not set up PKCE challenge: %s", err)
		}
		authFnInt, err := authFn(conf, state, challenge)
		if err != nil {
			return nil, fmt.Errorf("could not receive auth code: %s", err)
		}
		authCode, err := authFnInt()
		if err != nil {
			return nil, fmt.Errorf("could not receive auth code: %s", err)
		}
		authVerifyVal := oauth2.SetAuthURLParam("code_verifier", verifier)
		token, err = conf.Exchange(oauth2.NoContext, authCode, authVerifyVal)
		if err != nil {
			return nil, fmt.Errorf("failed to exchange auth code for token: %s", err)
		}
	}

	return oauth2.NewClient(
		oauth2.NoContext,
		FileSource(tokenFile, token, conf),
	), nil
}

func NewRefreshTokenClient(clientId, clientSecret, refreshToken string) *http.Client {
	conf := getConfig(clientId, clientSecret)

	token := &oauth2.Token{
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		Expiry:       time.Now(),
	}

	return oauth2.NewClient(
		oauth2.NoContext,
		conf.TokenSource(oauth2.NoContext, token),
	)
}

func NewAccessTokenClient(clientId, clientSecret, accessToken string) *http.Client {
	conf := getConfig(clientId, clientSecret)

	token := &oauth2.Token{
		TokenType:   "Bearer",
		AccessToken: accessToken,
	}

	return oauth2.NewClient(
		oauth2.NoContext,
		conf.TokenSource(oauth2.NoContext, token),
	)
}

func NewServiceAccountClient(serviceAccountFile string) (*http.Client, error) {
	content, exists, err := ReadFile(serviceAccountFile)
	if !exists {
		return nil, fmt.Errorf("Service account filename %q not found", serviceAccountFile)
	}

	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(content, "https://www.googleapis.com/auth/drive")
	if err != nil {
		return nil, err
	}
	return conf.Client(oauth2.NoContext), nil
}

func getConfig(clientId, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

func makeState() (string, error) {
	return makeString(12)
}

func makeCodeChallenge() (string, string, error) {
	verifier, err := makeString(48)
	if err != nil {
		return "", "", err
	}

	hasher := sha256.New()
	_, err = hasher.Write([]byte(verifier))
	if err != nil {
		return "", "", err
	}

	hash := hasher.Sum(nil)
	challenge := base64.RawURLEncoding.EncodeToString(hash)

	return verifier, challenge, nil
}

func makeString(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}
