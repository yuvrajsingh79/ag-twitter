package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-agauth/twitter/campaigns"
	twitterOAuth1 "github.com/go-agauth/twitter/oauth1"
	"github.com/go-agauth/twitter/sessions"
	"github.com/go-agauth/twitter/users"
)

const (
	sessionName     = "twitter-oauth-session"
	sessionSecret   = "twitter cookie signing secret"
	sessionUserKey  = "twitterID"
	sessionUsername = "twitterUsername"
	sessionAdCreds  = "twitterCredentials"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// Config configures the main ServeMux.
type Config struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
}

var adCreds = &campaigns.AdCreds{
	ConsumerKey:      "",
	ConsumerSecret:   "",
	OAuthToken:       "",
	OAuthTokenSecret: "",
	OAuthVerifier:    "",
}

// New returns a new ServeMux with app routes.
func New(config *Config) *http.ServeMux {
	mux := http.NewServeMux()

	//authentication begins...
	mux.HandleFunc("/", profileHandler)
	mux.HandleFunc("/logout", logoutHandler)
	// 1. Register Twitter login and callback handlers
	oauth1Config := &twitterOAuth1.Config{
		ConsumerKey:    config.TwitterConsumerKey,
		ConsumerSecret: config.TwitterConsumerSecret,
		CallbackURL:    "http://localhost:8080/twitter/callback",
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	mux.Handle("/twitter/login", twitterOAuth1.LoginController(oauth1Config, nil))
	mux.Handle("/twitter/callback", twitterOAuth1.CallbackController(oauth1Config, issueSession(oauth1Config), nil))
	//...authentication process completed

	//ads api route handlers
	fmt.Println("-----------------------****************************--------------------------")
	// adsConf := adsInit(oauth1Config)
	// fmt.Println(adsConf)
	// mux.Handle("/twitter/ads/accounts")

	return mux
}

// issueSession issues a cookie session after successful Twitter login
func issueSession(config *twitterOAuth1.Config) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := users.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		accessToken, accessSecret, err := twitterOAuth1.AccessTokenFromContext(ctx)
		_, verifier, err := twitterOAuth1.ParseAuthorizationCallback(req)
		adCreds = &campaigns.AdCreds{
			ConsumerKey:      config.ConsumerKey,
			ConsumerSecret:   config.ConsumerSecret,
			OAuthToken:       accessToken,
			OAuthTokenSecret: accessSecret,
			OAuthVerifier:    verifier,
		}

		fmt.Println("-----------------------****************************--------------------------")
		fmt.Println(adCreds)

		// 2. Implement a success handler to issue some form of session
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = twitterUser.ID
		session.Values[sessionUsername] = twitterUser.ScreenName
		session.Values[sessionAdCreds] = adCreds
		session.Save(w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// profileHandler shows a personal profile or a login button.
func profileHandler(w http.ResponseWriter, req *http.Request) {
	session, err := sessionStore.Get(req, sessionName)
	if err != nil {
		page, _ := ioutil.ReadFile("login.html")
		fmt.Fprintf(w, string(page))
		return
	}

	// authenticated profile
	fmt.Fprintf(w, `<p>You are logged in %s!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`, session.Values[sessionUsername])
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		sessionStore.Destroy(w, sessionName)
	}
	http.Redirect(w, req, "/", http.StatusFound)
}

// main creates and starts a Server listening.
func main() {
	const address = "localhost:8080"
	// read credentials from environment variables if available
	config := &Config{
		TwitterConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		TwitterConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	// allow consumer credential flags to override config fields
	consumerKey := flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flag.String("consumer-secret", "", "Twitter Consumer Secret")
	// consumerKey := "eaGx1xxDSGm3Mvng42VvRuJRT"
	// consumerSecret := "YDV0B453ckaWU5yyYBBHIG0SzZMoeRVqaEgGjB0ZRVYWtalJwI"
	flag.Parse()

	if *consumerKey != "" {
		config.TwitterConsumerKey = *consumerKey
	}

	if *consumerSecret != "" {
		config.TwitterConsumerSecret = *consumerSecret
	}

	if config.TwitterConsumerKey == "" {
		log.Fatal("Missing Twitter Consumer Key")
	}

	if config.TwitterConsumerSecret == "" {
		log.Fatal("Missing Twitter Consumer Secret")
	}

	log.Printf("Starting Server listening on %s\n", address)
	err := http.ListenAndServe(address, New(config))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
