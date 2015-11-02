package brightbox

import (
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
	CurrentToken *oauth2.Token
}

// Authenticate the details and return a client
// Region must be in regionURL map.
func (authd *authdetails) authenticatedClient() (*brightbox.Client, error) {
	switch {
	case authd.CurrentToken != nil:
		return authd.tokenisedAuth()
	case authd.UserName != "" || authd.password != "":
		return authd.tokenisedAuth()
	default:
		return authd.apiClientAuth()
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
	if authd.CurrentToken == nil {
		token, err := conf.PasswordCredentialsToken(oauth2.NoContext, authd.UserName, authd.password)
		if err != nil {
			return nil, err
		}
		authd.CurrentToken = token
	}
	oauth_connection := conf.Client(oauth2.NoContext, authd.CurrentToken)
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
