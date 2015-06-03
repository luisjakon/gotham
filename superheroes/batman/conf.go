package batman

import (
	"github.com/luisjakon/gotham"
	"net/http"
)

// Convenience utilities for batman initialization
type Conf struct {
	SecretKey       []byte
	GetProviderName func(r *http.Request) string
}

func Init(c Conf, provs ...gotham.Provider) error {

	GetProviderName = c.GetProviderName

	if c.SecretKey != nil {
		SecretKey(c.SecretKey)
	}

	if provs != nil {
		Protect(provs...)
	}

	return nil
}
