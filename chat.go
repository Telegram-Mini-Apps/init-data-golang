package twa

const (
	ChatTypeGroup      ChatType = "group"
	ChatTypeSupergroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)

// ChatType describes type of chat.
type ChatType = string

// Chat describes chat information:
// https://core.telegram.org/bots/webapps#webappchat
type Chat struct {
	// Unique identifier for this chat.
	Id int64 `json:"id"`

	// Optional. URL of the chat’s photo. The photo can be in .jpeg or .svg
	// formats. Only returned for Web Apps launched from the attachment menu.
	PhotoUrl string `json:"photo_url"`

	// Type of chat.
	Type ChatType `json:"type"`

	// Title of the chat.
	Title string `json:"title"`

	// Optional. Username of the chat.
	Username string `json:"username"`
}
