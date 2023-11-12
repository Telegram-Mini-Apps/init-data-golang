package initdata

const (
	ChatTypeSender     ChatType = "sender"
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeSupergroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)

// ChatType describes type of chat.
type ChatType string

// Known returns true if current chat type is known.
func (c ChatType) Known() bool {
	switch c {
	case ChatTypeSender,
		ChatTypePrivate,
		ChatTypeGroup,
		ChatTypeSupergroup,
		ChatTypeChannel:
		return true
	default:
		return false
	}
}

// Chat describes chat information:
// https://docs.telegram-mini-apps.com/launch-parameters/init-data#chat
type Chat struct {
	// Unique identifier for this chat.
	ID int64 `json:"id"`

	// Type of chat.
	Type ChatType `json:"type"`

	// Title of the chat.
	Title string `json:"title"`

	// Optional. URL of the chat’s photo. The photo can be in .jpeg or .svg
	// formats. Only returned for Web Apps launched from the attachment menu.
	PhotoURL string `json:"photo_url"`

	// Optional. Username of the chat.
	Username string `json:"username"`
}
