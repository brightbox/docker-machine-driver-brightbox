package brightbox

import (
	"github.com/brightbox/gobrightbox"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	defaultRegion = "gb1"
)

var regionURL = map[string]string{
	defaultRegion: brightbox.RegionGB1,
}

var infrastructureScope = []string{"infrastructure"}

type authdetails struct {
	APIClient string
	apiSecret string
	UserName  string
	password  string
	Account   string
	Region    string
}

// Authenticate the details and return a client
// Region must be in regionURL map.
func (authd *authdetails) authenticatedClient() (*brightbox.Client, error) {
	switch {
	case authd.UserName != "" || authd.password != "":
		return authd.passwordAuth()
	default:
		return authd.apiClientAuth()
	}
}

func (authd *authdetails) tokenURL() string {
	return authd.apiURL() + "/token"
}

func (authd *authdetails) apiURL() string {
	return regionURL[authd.Region]
}

func (authd *authdetails) passwordAuth() (*brightbox.Client, error) {
	conf := oauth2.Config{
		ClientID:     authd.APIClient,
		ClientSecret: authd.apiSecret,
		Scopes:       infrastructureScope,
		Endpoint: oauth2.Endpoint{
			TokenURL: authd.tokenURL(),
		},
	}
	token, err := conf.PasswordCredentialsToken(oauth2.NoContext, authd.UserName, authd.password)
	if err != nil {
		return nil, err
	}
	oauth_connection := conf.Client(oauth2.NoContext, token)
	return brightbox.NewClient(authd.apiURL(), authd.Account, oauth_connection)
}

func (authd *authdetails) apiClientAuth() (*brightbox.Client, error) {
	conf := clientcredentials.Config{
		ClientID:     authd.APIClient,
		ClientSecret: authd.apiSecret,
		Scopes:       infrastructureScope,
		TokenURL:     authd.tokenURL(),
	}
	oauth_connection := conf.Client(oauth2.NoContext)
	return brightbox.NewClient(authd.apiURL(), authd.Account, oauth_connection)
}
