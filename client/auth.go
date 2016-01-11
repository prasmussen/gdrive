package client

import (
    "net/http"
    "golang.org/x/oauth2"
    "go4.org/oauthutil"
)

type authCodeFn func(string) func() string

func NewOauthClient(clientId, clientSecret, cacheFile string, authFn authCodeFn) *http.Client {
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

    authUrl := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

    tokenSource := oauthutil.TokenSource{
        Config: conf,
        CacheFile: cacheFile,
        AuthCode: authFn(authUrl),
    }

    return oauth2.NewClient(oauth2.NoContext, tokenSource)
}
