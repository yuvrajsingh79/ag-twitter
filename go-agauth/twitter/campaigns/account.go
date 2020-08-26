package campaigns

import (
	"log"
	"net/http"

	"github.com/go-agauth/twitter/sessions"
)

const (
	sessionName      = "twitter-oauth-session"
	sessionSecret    = "twitter cookie signing secret"
	sessionConKey    = "ConsumerKey"
	sessionConSecret = "ConsumerSecret"
	sessionOToken    = "OAuthToken"
	sessionOSecret   = "OAuthSecret"
	sessionVerifier  = "OAuthVerifier"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

//AddNewAccount creates a new ad account
func AddNewAccount() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		// ctx := req.Context()
		session, err := sessionStore.Get(req, sessionName)
		if err != nil {
			log.Println("err")
		}
		log.Println(session.Values[sessionConKey], " ", session.Values[sessionConSecret], " ", session.Values[sessionOToken], " ", session.Values[sessionOSecret], " ", session.Values[sessionVerifier])
	}
	return http.HandlerFunc(fn)
}
