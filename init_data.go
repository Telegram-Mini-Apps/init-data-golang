package twa

import "time"

// InitData describes parsed initial data sent from TWA application. You can
// find specification for all the parameters in the official documentation:
// https://core.telegram.org/bots/webapps#webappinitdata
type InitData struct {
	// Init data generation date.
	AuthDateRaw int `json:"auth_date"`

	// Optional. Time in seconds, after which a message can be sent via the
	// `answerWebAppQuery` method.
	//
	// See: https://core.telegram.org/bots/api#answerwebappquery
	CanSendAfterRaw int `json:"can_send_after"`

	// An object containing data about the chat where the bot was
	// launched via the attachment menu. Returned for supergroups, channels
	// and group chats – only for Web Apps launched via the attachment menu.
	Chat *Chat `json:"chat"`

	// A hash of all passed parameters, which the bot server can use to
	// check their validity.
	//
	// See: https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app
	Hash string `json:"hash"`

	// A unique identifier for the Web App session, required for sending
	// messages via the answerWebAppQuery method.
	//
	// See: https://core.telegram.org/bots/api#answerwebappquery
	QueryID string `json:"query_id"`

	// An object containing data about the chat partner of the current user in
	// the chat where the bot was launched via the attachment menu.
	// Returned only for private chats and only for Web Apps launched
	// via the attachment menu.
	Receiver *User `json:"receiver"`

	// Optional. The value of the `startattach` parameter, passed via link. Only
	// returned for Web Apps when launched from the attachment menu via link.
	StartParam string `json:"start_param"`

	// An object containing data about the current user.
	User *User `json:"user"`
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
