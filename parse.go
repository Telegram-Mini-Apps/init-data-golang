package initdata

import (
	"encoding/json"
	"fmt"
	"net/url"
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
func Parse(initData string) (InitData, error) {
	// Parse passed init data as query string.
	q, err := url.ParseQuery(initData)
	if err != nil {
		return InitData{}, ErrUnexpectedFormat
	}

	// According to documentation, we could only meet such types as int64,
	// string, or another object. So, we create
	pairs := make([]string, 0, len(q))
	for k, v := range q {
		// Derive real value. We know that there can not be any arrays and value
		// can be the only one.
		val := v[0]
		valFormat := "%q:%q"

		// If passed value is valid in the context of JSON, it means, we could
		// insert this value without formatting.
		if isString := _stringProps[k]; !isString && json.Valid([]byte(val)) {
			valFormat = "%q:%s"
		}

		pairs = append(pairs, fmt.Sprintf(valFormat, k, val))
	}

	// Unmarshal JSON to our custom structure.
	var d InitData
	jStr := fmt.Sprintf("{%s}", strings.Join(pairs, ","))
	if err := json.Unmarshal([]byte(jStr), &d); err != nil {
		return InitData{}, ErrUnexpectedFormat
	}
	return d, nil
}
