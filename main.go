package twa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
// expIn - maximum init data lifetime. It is strongly recommended using this
// parameter. In case, zero duration is less than 0, function does not check if
// parameters are expired.
func Validate(initData, token string, expIn time.Duration) error {
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

	// Compute sign.
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(token))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(strings.Join(pairs, "\n")))

	// In case, our sign is not equal to found one, we should throw an error.
	if hex.EncodeToString(impHmac.Sum(nil)) != hash {
		return ErrSignInvalid
	}
	return nil
}

// Parse converts passed init data presented as query string to InitData
// object.
func Parse(initData string) (*InitData, error) {
	// Parse passed init data as query string.
	q, err := url.ParseQuery(initData)
	if err != nil {
		return nil, ErrUnexpectedFormat
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
		if json.Valid([]byte(val)) {
			valFormat = "%q:%s"
		}
		pairs = append(pairs, fmt.Sprintf(valFormat, k, val))
	}

	// Unmarshal JSON to our custom structure.
	var d InitData
	jStr := fmt.Sprintf("{%s}", strings.Join(pairs, ","))
	if err := json.Unmarshal([]byte(jStr), &d); err != nil {
		return nil, ErrUnexpectedFormat
	}
	return &d, nil
}
