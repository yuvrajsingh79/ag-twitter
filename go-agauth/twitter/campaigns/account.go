package campaigns

import (
	"log"
	"net/http"

	"github.com/go-agauth/twitter/sessions"
)

const (
	sessionAdCreds = "twitterCredentials"
	sessionName    = "twitter-oauth-session"
	sessionSecret  = "twitter cookie signing secret"
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
		log.Println(session.Values[sessionAdCreds])
	}
	return http.HandlerFunc(fn)
}
