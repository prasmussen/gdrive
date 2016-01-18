package auth

import (
    "net/http"
    "golang.org/x/oauth2"
)

type authCodeFn func(string) func() string

func NewOauthClient(clientId, clientSecret, tokenFile string, authFn authCodeFn) (*http.Client, error) {
    conf := &oauth2.Config{
        ClientID:     clientId,
        ClientSecret: clientSecret,
        Scopes:       []string{"https://www.googleapis.com/auth/drive"},
        RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
        Endpoint: oauth2.Endpoint{
            AuthURL:  "https://accounts.google.com/o/oauth2/auth",
            TokenURL: "https://accounts.google.com/o/oauth2/token",
        },
    }

    // Read cached token
    token, exists, err := ReadToken(tokenFile)
    if err != nil {
        return nil, err
    }

    // Require auth code if token file does not exist
    // or refresh token is missing
    if !exists || token.RefreshToken == "" {
        authUrl := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
        authCode := authFn(authUrl)()
        token, err = conf.Exchange(oauth2.NoContext, authCode)
    }

    return oauth2.NewClient(
        oauth2.NoContext,
        FileSource(tokenFile, token, conf),
    ), nil
}
