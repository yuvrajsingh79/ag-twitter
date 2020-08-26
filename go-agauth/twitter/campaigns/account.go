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

//Account model for Twitter Ads-Account
type Account struct {
	AccountID      string `json:"account_ids"`
	Count          int    `json:"count"`
	Cursor         string `json:"cursor"`
	QueryScope     string `json:"q"`
	SortBy         string `json:"sort_by"`
	WithDeleted    bool   `json:"with_deleted"`
	WithTotalCount bool   `json:"with_total_count"`
}

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

//AddNewAccount creates a new ad account
func AddNewAccount() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		session, err := sessionStore.Get(req, sessionName)
		if err != nil {
			log.Println("err")
		}
		log.Println(session.Values[sessionConKey], " ", session.Values[sessionConSecret], " ", session.Values[sessionOToken], " ", session.Values[sessionOSecret], " ", session.Values[sessionVerifier])

		//TODO --
	}
	return http.HandlerFunc(fn)
}

//NewAdAccount Create an ads account
func NewAdAccount() {

}

//NewSandboxAccount Create an ads account in the Sandbox environment
func NewSandboxAccount() {

}
