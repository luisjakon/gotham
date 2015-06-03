package gplus

import (
	"net/http"
	"net/url"

	"github.com/luisjakon/gotham"
	"golang.org/x/oauth2"
)

func New(clientId, clientSecret, callbackUrl string, scopes ...string) gotham.Provider {
	return &gplus{
		gotham.NewProvider("gplus", &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			RedirectURL:  callbackUrl,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://accounts.google.com/o/oauth2/token",
			},
		}),
		"https://www.googleapis.com/oauth2/v2/userinfo",
	}
}

// Google+ -  gotham.AuthProvider
type gplus struct {
	gotham.Provider
	ProfileURL string
}

func (p gplus) FetchUserData(t interface{}) (*gotham.UserData, error) {
	tokn := t.(*oauth2.Token)
	user := &gotham.UserData{AccessToken: tokn.AccessToken}

	res, err := http.Get(p.ProfileURL + "&access_token=" + url.QueryEscape(tokn.AccessToken))
	if err != nil {
		return user, err
	}

	return gotham.GetRawUserData(res, user)
}
