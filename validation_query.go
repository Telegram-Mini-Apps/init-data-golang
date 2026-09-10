package initdata

import (
	"fmt"
	"net/url"
)

func parseValidationQuery(initData string) (url.Values, error) {
	q, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("parse init data as query: %w: %w", err, ErrUnexpectedFormat)
	}

	for key, values := range q {
		if len(values) != 1 {
			return nil, fmt.Errorf(
				"parameter %q appears %d times: %w",
				key,
				len(values),
				ErrUnexpectedFormat,
			)
		}
	}
	return q, nil
}
