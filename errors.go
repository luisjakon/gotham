package gotham

import "errors"

var (
	ErrInvalidAuthState       = errors.New("Invalid Auth state.")
	ErrInvalidStateSignatures = errors.New("Invalid Auth state.")
	ErrInvalidToken           = errors.New("Invalid Token.")
	ErrExpiredStateCookie     = errors.New("Auth window expired.")
	ErrAuthProvider           = errors.New("Auth provider error.")
	ErrUnknownAuthProvider    = errors.New("Unknown auth provider error.")
	ErrNotImplemented         = errors.New("Method not implemented.")
)
