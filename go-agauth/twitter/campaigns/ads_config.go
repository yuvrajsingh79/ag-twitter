package campaigns

import (
	"fmt"
	"net/http"
)

//AdCreds hols the params of oauth1 credentials
type AdCreds struct {
	ConsumerKey      string `json:"consumer_key"`
	ConsumerSecret   string `json:"consumer_secret"`
	OAuthToken       string `json:"oauth-token"`
	OAuthTokenSecret string `json:"oauth-token_secret"`
	OAuthVerifier    string `json:"oauth_verifier"`
}

//GetAccountDetails fetches
func GetAccountDetails(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hi")
}
