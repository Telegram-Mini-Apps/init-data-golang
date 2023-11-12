package initdata

// User describes user information:
// https://docs.telegram-mini-apps.com/launch-parameters/init-data#user
type User struct {
	// Optional. True, if this user added the bot to the attachment menu.
	AddedToAttachmentMenu bool `json:"added_to_attachment_menu"`

	// Optional. True, if this user allowed the bot to message them.
	AllowsWriteToPm bool `json:"allows_write_to_pm"`

	// First name of the user or bot.
	FirstName string `json:"first_name"`

	// A unique identifier for the user or bot.
	ID int64 `json:"id"`

	// Optional. True, if this user is a bot. Returned in the `receiver` field
	// only.
	IsBot bool `json:"is_bot"`

	// Optional. True, if this user is a Telegram Premium user.
	IsPremium bool `json:"is_premium"`

	// Optional. Last name of the user or bot.
	LastName string `json:"last_name"`

	// Optional. Username of the user or bot.
	Username string `json:"username"`

	// Optional. IETF language tag of the user's language. Returns in user
	// field only.
	// https://en.wikipedia.org/wiki/IETF_language_tag
	LanguageCode string `json:"language_code"`

	// Optional. URL of the user’s profile photo. The photo can be in .jpeg or
	// .svg formats. Only returned for Web Apps launched from the
	// attachment menu.
	PhotoURL string `json:"photo_url"`
}
