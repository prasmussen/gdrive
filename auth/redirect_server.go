package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
)

type authorize struct{ authUrl string }
type redirect struct {
	done  chan string
	bad   chan bool
	state string
}

func (a authorize) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Location", a.authUrl)
	w.WriteHeader(302)
	fmt.Fprintln(w, "<html><head>")
	fmt.Fprintln(w, "<title>Redirect to authentication server</title>")
	fmt.Fprintln(w, "</head><body>")
	fmt.Fprintf(w, "Click <a href=\"%s\">here</a> to authorize gdrive to use Google Drive\n",
		a.authUrl)
	fmt.Fprintln(w, "</body></html>")
}

func (r redirect) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Printf("Could not parse form on /redirect: %s\n", err)
		w.WriteHeader(400)
		fmt.Fprintln(w, "<html><head>")
		fmt.Fprintln(w, "<title>Bad request</title>")
		fmt.Fprintln(w, "</head><body>")
		fmt.Fprintln(w, "Bad request: Missing authentication response")
		fmt.Fprintln(w, "</body></html>")
		return
	}
	if req.Form.Has("error") {
		fmt.Printf("authentication failed, server response is %s\n", req.Form.Get("error"))
		r.bad <- true
		fmt.Fprintln(w, "<html><head>")
		fmt.Fprintln(w, "<title>Google Drive authentication failed</title>")
		fmt.Fprintln(w, "</head><body>")
		fmt.Fprintf(w, "Authentication failed or refused: %s\n", req.Form.Get("error"))
		fmt.Fprintln(w, "</body></html>")
		return
	}

	if !req.Form.Has("code") || !req.Form.Has("state") {
		fmt.Println("redirect request is missing parameters")
		w.WriteHeader(400)
		fmt.Fprintln(w, "<html><head>")
		fmt.Fprintln(w, "<title>Bad request</title>")
		fmt.Fprintln(w, "</head><body>")
		fmt.Fprintln(w, "Bad request: response is missing the code or state parameters")
		fmt.Fprintln(w, "</body></html>")
		return
	}

	code := req.Form.Get("code")
	state := req.Form.Get("state")
	if state != r.state {
		fmt.Printf("Redirect state mismatch: %s vs %s", state, r.state)
		w.WriteHeader(400)
		fmt.Fprintln(w, "<html><head>")
		fmt.Fprintln(w, "<title>Bad request</title>")
		fmt.Fprintln(w, "</head><body>")
		fmt.Fprintln(w, "Bad request: response state mismatch")
		fmt.Fprintln(w, "</body></html>")
		return
	}
	fmt.Fprintln(w, "<html><head>")
	fmt.Fprintln(w, "<title>Authentication response received</title>")
	fmt.Fprintln(w, "</head><body>")
	fmt.Fprintln(w, "Authentication response has been received. Check the terminal where gdrive is running")
	fmt.Fprintln(w, "</body></html>")

	r.done <- code
}

func AuthCodeHTTP(conf *oauth2.Config, state, challenge string) (func() (string, error), error) {

	authChallengeMeth := oauth2.SetAuthURLParam("code_challenge_method", "S256")
	authChallengeVal := oauth2.SetAuthURLParam("code_challenge", challenge)

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	hostPort := ln.Addr().String()
	_, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	srv := &http.Server{Handler: mux}

	go func() {
		err := srv.Serve(ln)
		if err != http.ErrServerClosed {
			fmt.Printf("Cannot start http server: %s", err)
			os.Exit(1)
		}
	}()
	myconf := conf
	myconf.RedirectURL = fmt.Sprintf("http://127.0.0.1:%s/callback", port)

	authUrl := myconf.AuthCodeURL(state, oauth2.AccessTypeOffline, authChallengeMeth, authChallengeVal)
	authorizer := authorize{authUrl: authUrl}
	mux.Handle("/authorize", authorizer)
	callback := redirect{state: state,
		done: make(chan string, 1),
		bad:  make(chan bool, 1),
	}
	mux.Handle("/callback", callback)

	return func() (string, error) {
		var code string
		var err error
		fmt.Println("Authentication needed")
		fmt.Println("Go to the following url in your browser:")
		fmt.Printf("http://127.0.0.1:%s/authorize\n\n", port)
		fmt.Println("Waiting for authentication response")

		select {
		case <-callback.bad:
			err = fmt.Errorf("authentication did not complete successfully")
			code = ""
		case code = <-callback.done:
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()

		if stoperr := srv.Shutdown(ctx); stoperr != nil {
			fmt.Printf("Server Shutdown Failed:%+v\n", stoperr)
		}
		return code, err
	}, nil
}
