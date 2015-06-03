package gotham

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	root       = "/"            // auth state valid cookie path
	stateparam = "_oauth_state" // auth state cookiestore key name
	statelen   = 16             // auth state param length in bytes
	MaxAge     = 60             // auth timeout in seconds
)

var (
	NewState      = newState
	SaveState     = saveState
	GetSavedState = getSavedState
	ValidateState = validateState
	ClearState    = clearState
)

func newState(w http.ResponseWriter, r *http.Request) (state string) {
	state = generateRandomKey(statelen)
	SaveState(w, state)
	return
}

func saveState(w http.ResponseWriter, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:   stateparam,
		Value:  timestamp(state),
		Path:   root,
		MaxAge: MaxAge,
	})
}

func clearState(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    stateparam,
		Value:   root,
		Path:    root,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	})
}

func getSavedState(r *http.Request) string {
	c, err := r.Cookie(stateparam)
	if err != nil {
		panic(err.Error())
	}

	vals := strings.Split(c.Value, "::")
	timestamp, timesig, savedstate := vals[0], vals[1], vals[2]

	if !verify(timesig, timestamp) {
		return ErrInvalidAuthState.Error()
	}

	if isExpired(timestamp) {
		return ErrExpiredStateCookie.Error()
	}

	return savedstate
}

func validateState(r *http.Request) error {
	savedstate := GetSavedState(r)
	receivedstate := GetProviderState(r)

	if err := recover(); err != nil {
		return err.(error)
	}

	if !verify(savedstate, receivedstate) {
		return ErrInvalidStateSignatures
	}

	return nil
}

func isExpired(stamp string) bool {
	i, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		return false
	}
	expires := time.Unix(i, 0)
	return expires.Add(MaxAge * time.Second).Before(time.Now())
}
