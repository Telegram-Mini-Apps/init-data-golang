package initdata

import (
	"time"
)

// InitData contains init data.
// https://docs.telegram-mini-apps.com/launch-parameters/init-data#parameters-list
type InitData struct {
	// The date the initialization data was created. Is a number representing a
	// Unix timestamp.
	AuthDateRaw int `json:"auth_date"`

	// Optional. The number of seconds after which a message can be sent via
	// the method answerWebAppQuery.
	// https://core.telegram.org/bots/api#answerwebappquery
	CanSendAfterRaw int `json:"can_send_after"`

	// Optional. An object containing information about the chat with the bot in
	// which the Mini Apps was launched. It is returned only for Mini Apps
	// opened through the attachment menu.
	Chat Chat `json:"chat"`

	// Optional. The type of chat from which the Mini Apps was opened.
	// Returned only for applications opened by direct link.
	ChatType ChatType `json:"chat_type"`

	// Optional. A global identifier indicating the chat from which the Mini
	// Apps was opened. Returned only for applications opened by direct link.
	ChatInstance int64 `json:"chat_instance"`

	// Initialization data signature.
	// https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app
	Hash string `json:"hash"`

	// Optional. The unique session ID of the Mini App. Used in the process of
	// sending a message via the method answerWebAppQuery.
	// https://core.telegram.org/bots/api#answerwebappquery
	QueryID string `json:"query_id"`

	// Optional. An object containing data about the chat partner of the current
	// user in the chat where the bot was launched via the attachment menu.
	// Returned only for private chats and only for Mini Apps launched via the
	// attachment menu.
	Receiver User `json:"receiver"`

	// Optional. The value of the startattach or startapp query parameter
	// specified in the link. It is returned only for Mini Apps opened through
	// the attachment menu.
	StartParam string `json:"start_param"`

	// Optional. An object containing information about the current user.
	User User `json:"user"`
}

// AuthDate returns AuthDateRaw as time.Time.
func (d *InitData) AuthDate() time.Time {
	return time.Unix(int64(d.AuthDateRaw), 0)
}

// CanSendAfter returns computed time which depends on CanSendAfterRaw and
// AuthDate. Originally, CanSendAfterRaw means time in seconds, after which
// `answerWebAppQuery` method can be called and that's why this value could
// be computed as time.
func (d *InitData) CanSendAfter() time.Time {
	return d.AuthDate().Add(time.Duration(d.CanSendAfterRaw) * time.Second)
}
