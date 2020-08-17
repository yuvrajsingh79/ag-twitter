package oauth1

// AuthenticateEndpoint is Twitter's OAuth 1 endpoint which uses the
// oauth/authenticate AuthorizeURL redirect. Logged in users who have granted
// access are immediately authenticated and redirected to the callback URL.
var AuthenticateEndpoint = Endpoint{
	RequestTokenURL: "https://api.twitter.com/oauth/request_token",
	AuthorizeURL:    "https://api.twitter.com/oauth/authenticate",
	AccessTokenURL:  "https://api.twitter.com/oauth/access_token",
}

// AuthorizeEndpoint is Twitter's OAuth 1 endpoint which uses the
// oauth/authorize AuthorizeURL redirect. Note that this requires users who
// have granted access previously, to re-grant access at AuthorizeURL.
// Prefer AuthenticateEndpoint over AuthorizeEndpoint if you are unsure.
var AuthorizeEndpoint = Endpoint{
	RequestTokenURL: "https://api.twitter.com/oauth/request_token",
	AuthorizeURL:    "https://api.twitter.com/oauth/authorize",
	AccessTokenURL:  "https://api.twitter.com/oauth/access_token",
}
