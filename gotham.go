package gotham

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

/**
 * Gotham - Lightweight Go Authentication Manager (https://github.com/luisjakon/gotham)
 */

type Providers map[string]Provider

var (
	Use                  = useProvider // legacy purposes only
	UseProvider          = useProvider
	GetProvider          = getProvider
	GetProviderName      = getProviderName
	GetProviderState     = getProviderState
	GetAuthorizationCode = getAuthCode
	GetRawUserData       = getRawUserData
)

var providers = Providers{}

func useProvider(provs ...Provider) Providers {
	for _, p := range provs {
		providers[p.Name()] = p
	}
	return providers
}

func getProvider(name string) Provider {
	return providers[name]
}

func getProviderName(r *http.Request) string {
	panic("undefined gotham.GetProviderName() function")
}

func getProviderState(r *http.Request) string {
	return r.URL.Query().Get("state")
}

func getAuthCode(r *http.Request) string {
	return r.URL.Query().Get("code")
}

func getRawUserData(res *http.Response, user *UserData) (*UserData, error) {
	if res == nil {
		return user, nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return user, err
	}

	return user, json.NewDecoder(bytes.NewReader(body)).Decode(&user.RawData)
}
