package initdata

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var (
	// List of properties which should always be interpreted as strings.
	_stringProps = map[string]bool{
		"start_param": true,
	}
)

// Parse converts passed init data presented as query string to InitData
// object.
func Parse(initData string) (*InitData, error) {
	// Parse passed init data as query string.
	q, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("parse init data as query: %w: %w", err, ErrUnexpectedFormat)
	}

	var d InitData
	for k, v := range q {
		target := d.fieldByName(k)
		if target == nil {
			continue
		}

		// Derive real value. We know that there can not be any arrays and value
		// can be the only one.
		if err := decodeValue(k, v[0], target); err != nil {
			return nil, fmt.Errorf("unmarshal init data: %w: %w", err, ErrUnexpectedFormat)
		}
	}
	return &d, nil
}

// fieldByName returns a pointer to the InitData field matching the passed
// property name. Names are matched case-insensitively, the same way
// encoding/json matches object keys.
func (d *InitData) fieldByName(name string) any {
	switch {
	case strings.EqualFold(name, "auth_date"):
		return &d.AuthDateRaw
	case strings.EqualFold(name, "can_send_after"):
		return &d.CanSendAfterRaw
	case strings.EqualFold(name, "chat"):
		return &d.Chat
	case strings.EqualFold(name, "chat_type"):
		return &d.ChatType
	case strings.EqualFold(name, "chat_instance"):
		return &d.ChatInstance
	case strings.EqualFold(name, "hash"):
		return &d.Hash
	case strings.EqualFold(name, "query_id"):
		return &d.QueryID
	case strings.EqualFold(name, "receiver"):
		return &d.Receiver
	case strings.EqualFold(name, "start_param"):
		return &d.StartParam
	case strings.EqualFold(name, "user"):
		return &d.User
	default:
		return nil
	}
}

// decodeValue decodes a single property value into target. According to
// documentation, we could only meet such types as int64, string, or another
// object. If the value is valid JSON, it is decoded as is, otherwise it is
// treated as a plain string.
func decodeValue(name, value string, target any) error {
	switch t := target.(type) {
	case *int:
		if n, err := strconv.Atoi(value); err == nil && isJSONInteger(value) {
			*t = n
			return nil
		}
	case *int64:
		if n, err := strconv.ParseInt(value, 10, 64); err == nil && isJSONInteger(value) {
			*t = n
			return nil
		}
	case **User, **Chat:
		// A value which is not valid JSON can not be decoded into an object
		// either way, so there is no need to validate it separately.
		return json.Unmarshal([]byte(value), target)
	}

	if !_stringProps[name] && json.Valid([]byte(value)) {
		return json.Unmarshal([]byte(value), target)
	}

	switch t := target.(type) {
	case *string:
		*t = value
	case *ChatType:
		*t = ChatType(value)
	default:
		return fmt.Errorf("property %q: cannot use string value as %T", name, target)
	}
	return nil
}

// isJSONInteger reports whether value is an integer written the way JSON
// allows it: an optional minus sign and digits without leading zeros.
func isJSONInteger(value string) bool {
	digits := strings.TrimPrefix(value, "-")
	if digits == "" || digits[0] == '+' {
		return false
	}
	return digits == "0" || digits[0] != '0'
}
