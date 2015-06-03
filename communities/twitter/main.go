package twitter

import (
	"net/http"
	"strings"

	"github.com/luisjakon/gotham"
	"github.com/mrjones/oauth"
)

func New(clientId, clientSecret, callbackUrl string) gotham.Provider {
	return &twitter{
		name: "twitter",
		consumer: oauth.NewConsumer(
			clientId,
			clientSecret,
			oauth.ServiceProvider{
				RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
				AuthorizeTokenUrl: "https://api.twitter.com/oauth/authenticate",
				AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
			}),
		callbackURL: callbackUrl,
		profileURL:  "https://api.twitter.com/1.1/account/verify_credentials.json",
	}
}

// Twitter Map - holds requestTokens for pending user auths
var tokens = map[string]*oauth.RequestToken{}

func Reset() {
	tokens = map[string]*oauth.RequestToken{}
}

// Twitter - gotham.AuthProvider
type twitter struct {
	name        string
	callbackURL string
	profileURL  string
	consumer    *oauth.Consumer
}

func (p *twitter) Name() string {
	return p.name
}

func (p *twitter) BeginUserAuth(w http.ResponseWriter, r *http.Request) {
	requestToken, url, err := p.consumer.GetRequestTokenAndUrl(p.callbackURL)
	if err != nil {
		panic(err)
	}

	saveRequestToken(w, r, requestToken)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (p *twitter) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	requestToken := getSavedRequestToken(w, r)

	verificationCode := r.URL.Query().Get("oauth_verifier")

	return p.consumer.AuthorizeToken(requestToken, verificationCode)
}

func (p *twitter) FetchUserData(t interface{}) (*gotham.UserData, error) {
	tokn := t.(*oauth.AccessToken)
	user := &gotham.UserData{AccessToken: tokn.Token, AccessTokenSecret: tokn.Secret}

	// fetch user data
	res, err := p.consumer.Get(
		p.profileURL,
		map[string]string{"include_entities": "false", "skip_status": "true"},
		tokn,
	)
	if err != nil {
		return user, err
	}

	return DecodeUserData(res, user)
}

// Helper Functions
var DecodeUserData = decodeUserData

func decodeUserData(r *http.Response, user *gotham.UserData) (*gotham.UserData, error) {
	_, err := gotham.GetRawUserData(r, user)
	defer recover()

	n := user.RawData["name"].(string)
	name := strings.Split(n, "")

	user.FirstName = name[0]
	user.LastName = name[1]
	user.NickName = user.RawData["screen_name"].(string)
	user.Description = user.RawData["description"].(string)
	user.AvatarURL = user.RawData["profile_image_url"].(string)
	user.UserID = user.RawData["id_str"].(string)
	user.Location = user.RawData["location"].(string)

	return user, err
}

func saveRequestToken(w http.ResponseWriter, r *http.Request, t *oauth.RequestToken) {
	state := gotham.NewState(w, r)
	gotham.SaveState(w, state)

	tokens[gotham.Sign(state)] = t // <- store token req with signed_state as key
}

func getSavedRequestToken(w http.ResponseWriter, r *http.Request) (t *oauth.RequestToken) {
	state := gotham.GetSavedState(r)
	defer gotham.ClearState(w)
	defer delete(tokens, state)

	t = tokens[state] // <- get stored token req with signed_state (from provider) as key
	return
}
