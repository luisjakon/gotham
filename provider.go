package gotham

import (
	"golang.org/x/oauth2"
	"net/http"
)

type Provider interface {
	Name() string
	BeginUserAuth(http.ResponseWriter, *http.Request)
	CompleteUserAuth(http.ResponseWriter, *http.Request) (token interface{}, err error) //  returns x.oauth or x.oauth2 token depending on impl
	FetchUserData(token interface{}) (*UserData, error)                                 // consumes x.oauth or x.oauth2 token depending on impl
}

type provider struct {
	*oauth2.Config
	name string
}

func NewProvider(name string, config *oauth2.Config) Provider {
	return &provider{
		Config: config,
		name:   name,
	}
}

func (p *provider) BeginUserAuth(w http.ResponseWriter, r *http.Request) {
	state := NewState(w, r)
	url := p.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (p *provider) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	defer ClearState(w) // remove saved auth_state when done
	if err := ValidateState(r); err != nil {
		return nil, err
	}
	return p.Exchange(oauth2.NoContext, GetAuthorizationCode(r))
}

func (p *provider) FetchUserData(t interface{}) (*UserData, error) {
	return nil, ErrNotImplemented
}

func (p *provider) Name() string {
	return p.name
}
