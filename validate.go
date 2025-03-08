package initdata

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Validate validates passed init data. This method expects initData to be
// passed in the exact raw format as it could be found
// in window.Telegram.WebApp.initData.
//
// Returns error if something is wrong with the passed init data. Nil otherwise.
//
// initData - init data passed from application;
// token - init data Telegram Bot issuer token;
// expIn - maximum init data lifetime. It is strongly recommended to use this
// parameter. In case, exp duration is less than or equal to 0, function does
// not check if parameters are expired.
func Validate(initData, token string, expIn time.Duration) error {
	// Parse passed init data as query string.
	q, err := url.ParseQuery(initData)
	if err != nil {
		return fmt.Errorf("parse init data as query: %w: %w", err, ErrUnexpectedFormat)
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
			i, err := strconv.ParseInt(v[0], 10, 64)
			if err != nil {
				return fmt.Errorf("parse auth_date to int64: %w: %w", err, ErrAuthDateInvalid)
			}

			authDate = time.Unix(i, 0)
		}
		// Append new pair.
		pairs = append(pairs, k+"="+v[0])
	}

	// Sign is always required.
	if hash == "" {
		return ErrSignMissing
	}

	// In case, expiration time is passed, we do additional parameters check.
	if expIn > 0 {
		// In case, auth date is zero, it means, we can not check if parameters
		// are expired.
		if authDate.IsZero() {
			return ErrAuthDateMissing
		}

		// Check if init data is expired.
		if authDate.Add(expIn).Before(time.Now()) {
			return ErrExpired
		}
	}

	// According to docs, we sort all the pairs in alphabetical order.
	sort.Strings(pairs)

	// In case, our sign is not equal to found one, we should throw an error.
	if sign(strings.Join(pairs, "\n"), token) != hash {
		return ErrSignInvalid
	}
	return nil
}
