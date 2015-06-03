# Gotham - Go Authentication Manager

Heavily inspired by Mark Bates' Goth (https://github.com/markbates/goth)

<b>Why Gotham?</b>
 - Primarily for me to practice Go and understand the OAuth flows for multiple providers and versions
 - To enable automatic auth state protection by default
 - To be able to control & detach the fetch-userdata phase from the normal auth flow as needed (e.g. for token only requests)
 - Also because Go is easy and fun...

<b>Differences + Additions to Goth/Gothic</b>
 - Storeless, Sessionless, minimalist approach (server-side)
 - Automatic protection and validation of oauth states
 - App-defined security keys & strengths
 - App-defined auth request timeout periods
 - App-defined authentication flows (e.g. auto-fetch userdata or fetch token only)
 - App-defined global + per-provider userdata readers/decoders

<b>Future Work:</b>
 - Add more providers
 - Add authorization bindings

# Example
Batman Begins (https://github.com/luisjakon/gotham/blob/master/superheroes/example.go)

# Contributing
Pull requests, contributions, issue(s) reporting and feedback are welcome and encouraged.
