package initdata

import (
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Validate validates passed init data. This method expects initData to be
// passed in the exact raw format as it could be found
// in window.Telegram.WebApp.initData. Returns true in case init data is
// signed correctly, and it is allowed to trust it.
//
// Current code is implementation of algorithmic code described in official
// docs:
// https://core.telegram.org/bots/webapps#validating-data-received-via-the-web-app
//
// initData - init data passed from application;
// token - TWA bot secret token which was used to create init data;
// deadline - maximum init data lifetime. It is strongly recommended to use this
// parameter. In case, exp duration is less than or equal to 0, function does
// not check if parameters are expired.
func Validate(initData, token string, deadline time.Time) error {
	// Parse passed init data as query string.
	q, err := url.ParseQuery(initData)
	if err != nil {
		return ErrUnexpectedFormat
	}

	var (
		// Init data creation time.
		authDate time.Time
		// Init data sign.
		hash string
		// All found key-value pairs.
		pairs = make([]string, 0, len(q))
	)

	// Iterate over all key-value pairs of parsed parameters.
	for k, v := range q {
		// Store found sign.
		if k == "hash" {
			hash = v[0]
			continue
		}
		if k == "auth_date" {
			if i, err := strconv.Atoi(v[0]); err == nil {
				authDate = time.Unix(int64(i), 0)
			}
		}
		// Append new pair.
		pairs = append(pairs, k+"="+v[0])
	}

	// Sign is always required.
	if hash == "" {
		return ErrSignMissing
	}

	if authDate.IsZero() {
		return ErrAuthDateMissing
	}

	if !deadline.IsZero() && authDate.After(deadline) {
		return ErrExpired
	}

	// According to docs, we sort all the pairs in alphabetical order.
	sort.Strings(pairs)

	// In case, our sign is not equal to found one, we should throw an error.
	if sign(strings.Join(pairs, "\n"), token) != hash {
		return ErrSignInvalid
	}
	return nil
}
