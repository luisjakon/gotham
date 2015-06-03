package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/bmizerany/pat"

	"github.com/luisjakon/gotham/communities/facebook"
	"github.com/luisjakon/gotham/communities/gplus"
	"github.com/luisjakon/gotham/communities/twitter"

	"github.com/luisjakon/gotham/superheroes/batman"
)

// Init
func init() {
	// Just in case
	os.Setenv("FACEBOOK_KEY", "your_fb_api_key")
	os.Setenv("FACEBOOK_SECRET", "your_fb_secret")
	os.Setenv("FACEBOOK_CALLBACK_URL", "http://localhost:3000/auth/facebook/callback")
	os.Setenv("TWITTER_KEY", "your_twitter_api_key")
	os.Setenv("TWITTER_SECRET", "your_twitter_secret")
	os.Setenv("TWITTER_CALLBACK_URL", "http://localhost:3000/auth/twitter/callback")
}

// Main
func main() {

	// Initialize batman
	batman.Init(batman.Conf{
		SecretKey: []byte("1234567890abcdef"), // enable auth req state protection against spoofing, forgeries & MITM attacks
		GetProviderName: func(r *http.Request) string { // enable custom provider id/name retrieval from the user auth request route/params
			return r.URL.Query().Get(":provider")
		},
	})

	// Register community auth providers
	batman.Protect(
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), os.Getenv("FACEBOOK_CALLBACK_URL"), "read_stream", "user_location"),
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), os.Getenv("TWITTER_CALLBACK_URL")),
		gplus.New(os.Getenv("GOOGLE+_KEY"), os.Getenv("GOOGLE+_SECRET"), os.Getenv("GOOGLE+_CALLBACK_URL")),
	)

	// Create router
	p := pat.New()

	// Display user auth options
	p.Get("/", wrap(home))

	// Handle user auth requests
	p.Get("/auth/:provider", batman.Begins)
	p.Get("/auth/:provider/callback", wrap(func(w http.ResponseWriter, r *http.Request) {
		provider, token, err := batman.Finalize(w, r)
		if err != nil {
			loginTemplate.Execute(w, err.Error())
			return
		}

		userdata, err := batman.FetchUserData(provider, token)
		if err != nil {
			loginTemplate.Execute(w, err.Error())
			return
		}

		profileTemplate.Execute(w, userdata)
	}))

	// Run batman
	http.ListenAndServe(":3000", p)
}

// Home
func home(w http.ResponseWriter, r *http.Request) {
	loginTemplate.Execute(w, "Batman welcomes you Gotham!")
}

// Wrap
func wrap(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(fn)
}

var loginTemplate = template.Must(template.New("").Parse(`
	<html>
	<head><link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap.min.css"></head>
	<body>
	<h1 class="text-center">{{.}}</h1>
	<h3 class="text-center">Please choose your gated community:</h2>
    <div class="row">
	<div class="col-md-8 col-md-offset-2">
	<div class="col-md-4">
	<a href="/auth/facebook" class="btn btn-lg btn-primary btn-block">Facebook</a>
	</div>
	<div class="col-md-4">
	<a href="/auth/twitter" class="btn btn-lg btn-info btn-block">Twitter</a>
	</div>
	<div class="col-md-4">
	<a href="/auth/gplus" class="btn btn-lg btn-danger btn-block">Google</a>
	</div>
	</div>
    </div>
	</body>
	</html>
	`))

var profileTemplate = template.Must(template.New("").Parse(`
	<html>
	<head><link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/css/bootstrap.min.css"></head>
	<body>
	<h1 class="text-center">Hi {{.NickName}}!</h1>
	<h3 class="text-center">Your Profile</h3>
	<div class="row">
	<div class="col-sm-6 col-md-8">
	<img src="{{.AvatarURL}}" alt="" class="img-rounded img-responsive text-center" />
	<h4>{{.FirstName}} {{.LastName}}</h4>
	<br />
	<small>{{.Description}}
    <i class="glyphicon glyphicon-user"></i>{{.NickName}} ({{.UserID}})
	<br />
	<i class="glyphicon glyphicon-envelope"></i>{{.Email}}
    <br />
    <i class="glyphicon glyphicon-globe"></i>{{.Location}}
    <br />
	<i class="glyphicon glyphicon-lock"></i>{{.AccessToken}}
	<br />
	<i class="glyphicon glyphicon-list-alt"></i>{{.RawData}}
	</div>
    </body></html>
	`))
