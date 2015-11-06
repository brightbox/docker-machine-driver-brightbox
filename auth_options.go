package brightbox

import (
	"os"

	"github.com/brightbox/gobrightbox"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var infrastructureScope = []string{"infrastructure"}

type authdetails struct {
	APIClient    string
	APISecret    string
	UserName     string
	password     string
	Account      string
	ApiURL       string
	currentToken *oauth2.Token
}

// Authenticate the details and return a client
// Region must be in regionURL map.
func (authd *authdetails) authenticatedClient() (*brightbox.Client, error) {
	authd.backfillPassword()
	switch {
	case authd.currentToken != nil:
		return authd.tokenisedAuth()
	case authd.UserName != "" || authd.password != "":
		return authd.tokenisedAuth()
	default:
		return authd.apiClientAuth()
	}
}

func (authd *authdetails) backfillPassword() {
	if authd.UserName != "" && authd.password == "" {
		authd.password = os.Getenv(passwordEnvVar)
	}
}

func (authd *authdetails) tokenURL() string {
	return authd.ApiURL + "/token"
}

func (authd *authdetails) tokenisedAuth() (*brightbox.Client, error) {
	conf := oauth2.Config{
		ClientID:     authd.APIClient,
		ClientSecret: authd.APISecret,
		Scopes:       infrastructureScope,
		Endpoint: oauth2.Endpoint{
			TokenURL: authd.tokenURL(),
		},
	}
	if authd.currentToken == nil {
		token, err := conf.PasswordCredentialsToken(oauth2.NoContext, authd.UserName, authd.password)
		if err != nil {
			return nil, err
		}
		authd.currentToken = token
	}
	oauth_connection := conf.Client(oauth2.NoContext, authd.currentToken)
	return brightbox.NewClient(authd.ApiURL, authd.Account, oauth_connection)
}

func (authd *authdetails) apiClientAuth() (*brightbox.Client, error) {
	conf := clientcredentials.Config{
		ClientID:     authd.APIClient,
		ClientSecret: authd.APISecret,
		Scopes:       infrastructureScope,
		TokenURL:     authd.tokenURL(),
	}
	oauth_connection := conf.Client(oauth2.NoContext)
	return brightbox.NewClient(authd.ApiURL, authd.Account, oauth_connection)
}
