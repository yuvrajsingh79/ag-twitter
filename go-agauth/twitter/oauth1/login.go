package oauth1

import (
	"errors"
	"net/http"

	"github.com/go-agauth/twitter/services"
	"github.com/go-agauth/twitter/users"
)

// Twitter login errors
var (
	ErrUnableToGetTwitterUser = errors.New("twitter: unable to get Twitter User")
)

// LoginController handles Twitter login requests by obtaining a request token and
// redirecting to the authorization URL.
func LoginController(config *Config, failure http.Handler) http.Handler {
	// oauth1.LoginHandler -> oauth1.AuthRedirectHandler
	success := AuthRedirectHandler(config, failure)
	return LoginHandler(config, success, failure)
}

// CallbackController handles Twitter callback requests by parsing the oauth token
// and verifier and adding the Twitter access token and User to the ctx. If
// authentication succeeds, handling delegates to the success handler,
// otherwise to the failure handler.
func CallbackController(config *Config, success, failure http.Handler) http.Handler {
	// oauth1.EmptyTempHandler -> oauth1.CallbackHandler -> TwitterHandler -> success
	success = twitterHandler(config, success, failure)
	success = CallbackHandler(config, success, failure)
	return EmptyTempHandler(success)
}

// twitterHandler is a http.Handler that gets the OAuth1 access token from
// the ctx and calls Twitter verify_credentials to get the corresponding User.
// If successful, the User is added to the ctx and the success handler is
// called. Otherwise, the failure handler is called.
func twitterHandler(config *Config, success, failure http.Handler) http.Handler {
	if failure == nil {
		failure = DefaultFailureHandler
	}
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		accessToken, accessSecret, err := AccessTokenFromContext(ctx)
		if err != nil {
			ctx = WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		httpClient := config.Client(ctx, NewToken(accessToken, accessSecret))
		twitterClient := services.NewClient(httpClient)
		accountVerifyParams := &services.AccountVerifyParams{
			IncludeEntities: services.Bool(false),
			SkipStatus:      services.Bool(true),
			IncludeEmail:    services.Bool(false),
		}
		user, resp, err := twitterClient.Accounts.VerifyCredentials(accountVerifyParams)
		err = validateResponse(user, resp, err)
		if err != nil {
			ctx = WithError(ctx, err)
			failure.ServeHTTP(w, req.WithContext(ctx))
			return
		}
		ctx = users.WithUser(ctx, user)
		success.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// validateResponse returns an error if the given Twitter user, raw
// http.Response, or error are unexpected. Returns nil if they are valid.
func validateResponse(user *users.User, resp *http.Response, err error) error {
	if err != nil || resp.StatusCode != http.StatusOK {
		return ErrUnableToGetTwitterUser
	}
	if user == nil || user.ID == 0 || user.IDStr == "" {
		return ErrUnableToGetTwitterUser
	}
	return nil
}
