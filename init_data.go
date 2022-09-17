package twa

import "time"

// InitData describes parsed initial data sent from TWA application. You can
// find specification for all the parameters in the official documentation:
// https://core.telegram.org/bots/webapps#webappinitdata
type InitData struct {
	QueryId         string `json:"query_id"`
	User            *User  `json:"user"`
	Receiver        *User  `json:"receiver"`
	Chat            *Chat  `json:"chat"`
	StartParam      string `json:"start_param"`
	CanSendAfterRaw int    `json:"can_send_after"`
	AuthDateRaw     int    `json:"auth_date"`
	Hash            string `json:"hash"`
}

// AuthDate returns AuthDateRaw as time.Time..
func (d *InitData) AuthDate() time.Time {
	return time.Unix(int64(d.AuthDateRaw), 0)
}

// CanSendAfter returns CanSendAfterRaw as time.Duration.
func (d *InitData) CanSendAfter() time.Duration {
	return time.Duration(d.CanSendAfterRaw) * time.Second
}
