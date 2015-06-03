package facebook

import (
	"net/http"
	"net/url"

	"github.com/luisjakon/gotham"
	"golang.org/x/oauth2"
)

func New(clientId, clientSecret, callbackUrl string, scopes ...string) gotham.Provider {
	return &facebook{
		gotham.NewProvider("facebook", &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			RedirectURL:  callbackUrl,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.facebook.com/dialog/oauth",
				TokenURL: "https://graph.facebook.com/oauth/access_token",
			},
		}),
		"https://graph.facebook.com/me?fields=email,first_name,last_name,link,bio,id,name,picture,location",
	}
}

// Facebook - gotham.AuthProvider
type facebook struct {
	gotham.Provider
	ProfileURL string
}

func (p facebook) FetchUserData(t interface{}) (*gotham.UserData, error) {
	tokn := t.(*oauth2.Token)
	user := &gotham.UserData{AccessToken: tokn.AccessToken}

	res, err := http.Get(p.ProfileURL + "&access_token=" + url.QueryEscape(tokn.AccessToken))
	if err != nil {
		return user, err
	}

	return DecodeUserData(res, user)
}

// Pluggable UserData Decoder
var DecodeUserData = decodeUserData

func decodeUserData(r *http.Response, user *gotham.UserData) (*gotham.UserData, error) {
	_, err := gotham.GetRawUserData(r, user)
	defer recover()

	user.FirstName = user.RawData["first_name"].(string)
	user.LastName = user.RawData["last_name"].(string)
	user.NickName = user.RawData["name"].(string)
	user.AvatarURL = user.RawData["picture"].(map[string]interface{})["data"].(map[string]interface{})["url"].(string)
	user.UserID = user.RawData["id"].(string)
	user.Email = user.RawData["email"].(string)

	return user, err
}
