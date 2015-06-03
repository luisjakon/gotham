package gotham

type UserData struct {
	RawData           map[string]interface{}
	Provider          string
	Email             string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	UserID            string
	AvatarURL         string
	Location          string
	Phone             string
	AccessToken       string
	AccessTokenSecret string
}
