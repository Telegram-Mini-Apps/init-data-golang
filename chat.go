package twa

// Chat describes chat information:
// https://core.telegram.org/bots/webapps#webappchat
type Chat struct {
	Id       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Username string `json:"username"`
	PhotoUrl string `json:"photo_url"`
}
