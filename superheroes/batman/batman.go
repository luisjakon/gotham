package batman

import (
	// "log"
	"net/http"

	"github.com/luisjakon/gotham"
	"golang.org/x/oauth2"
)

/**
 * Batman is to gotham as gothic is to goth (see https://github.com/markbates/goth)
 *
 * Batman Setup:
 *   STEP 1. - batman.GetProviderName - tells batman how to get provider id from an http request
 *   STEP 2. - batman.SecretKey(...) - initializes batman/gotham secret key
 *   STEP 3. - batman.Protect(...) - registers community auth providers
 *
 * Batman Workflow:
 *   STEP 1. - batman.Begin(...) - begins the auth exchange
 *   STEP 2. - batman.Finalize(...) - completes the auth exchange & return fresh oauth token
 *   STEP 3. - batman.FetchUserData(...) - retrieves user profile data from provider (OPTIONAL)
 */
var (
	Begins          = http.HandlerFunc(Begin) // Default batman http.HandlerFunc
	GetProviderName = gotham.GetProviderName  // Default batman.GetProviderName = func(*http.Request)string
)

func SecretKey(key []byte) {
	gotham.SecretKey = key
}

func Protect(provs ...gotham.Provider) {
	gotham.UseProvider(provs...)
}

func Begin(w http.ResponseWriter, r *http.Request) {
	provider := gotham.GetProvider(GetProviderName(r))
	if provider == nil {
		http.Redirect(w, r, "/auth/"+GetProviderName(r)+"/callback?error=unknown+provider", 302)
		return
	}
	provider.BeginUserAuth(w, r)
}

func Finalize(w http.ResponseWriter, r *http.Request) (id string, token interface{}, err error) {
	id = GetProviderName(r)

	provider := gotham.GetProvider(id)
	if provider == nil {
		token, err = &oauth2.Token{}, gotham.ErrUnknownAuthProvider
		return
	}

	token, err = provider.CompleteUserAuth(w, r)
	return
}

func FetchUserData(prov string, token interface{}) (*gotham.UserData, error) {
	provider := gotham.GetProvider(prov)
	if provider == nil {
		return &gotham.UserData{}, gotham.ErrUnknownAuthProvider
	}
	return provider.FetchUserData(token)
}
