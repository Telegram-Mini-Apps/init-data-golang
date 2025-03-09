package initdata

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	_telegramProdPublicKey, _ = hex.DecodeString("e7bf03a2fa4602af4580703d88dda5bb59f32ed8b02a56c187fe7d34caed242d")
	_telegramTestPublicKey, _ = hex.DecodeString("40055058a4ee38156a06562e52eece92a771bcd8346a8c4615cb7376eddf72ec")
)

// ValidateThirdPartyWithEnv validates passed init data assuming that it was signed by Telegram.
// This method expects initData to be passed in the exact raw format as it could be found
// in window.Telegram.WebApp.initData.
//
// Returns error if something is wrong with the passed init data. Nil otherwise.
//
// initData - init data passed from application;
// botID - init data Telegram Bot issuer identifier;
// expIn - maximum init data lifetime. It is strongly recommended to use this
// parameter. In case, exp duration is less than or equal to 0, function does
// not check if parameters are expired.
// isTest - true if the init data was issued in Telegram production environment;
func ValidateThirdPartyWithEnv(
	initData string,
	botID int64,
	expIn time.Duration,
	isTest bool,
) error {
	// Parse passed init data as query string.
	q, err := url.ParseQuery(initData)
	if err != nil {
		return fmt.Errorf("parse init data as query: %w: %w", err, ErrUnexpectedFormat)
	}

	var (
		// Init data creation time.
		authDate time.Time
		// Init data signature.
		signature []byte
		// All found key-value pairs.
		pairs = make([]string, 0, len(q))
	)

	// Iterate over all key-value pairs of parsed parameters.
	for k, v := range q {
		// hash is ignored during this type of validation.
		if k == "hash" {
			continue
		}

		// Store found signature.
		if k == "signature" {
			signature, _ = base64.URLEncoding.DecodeString(
				// Signature is base64-encoded. We append padding by ourselves as long Telegram's server
				// incorrectly creates a base64 string, but GoLang intolerant to this and requires strict
				// format compliance.
				v[0] + strings.Repeat("=", 4-len(v[0])%4),
			)
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

	// Signature is always required.
	if len(signature) == 0 {
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

	var publicKey []byte
	if isTest {
		publicKey = _telegramTestPublicKey
	} else {
		publicKey = _telegramProdPublicKey
	}

	if !ed25519.Verify(
		publicKey,
		[]byte(fmt.Sprintf("%d:WebAppData\n%s", botID, strings.Join(pairs, "\n"))),
		signature,
	) {
		return ErrSignInvalid
	}
	return nil
}

// ValidateThirdParty performs validation described in the Validate3rdWithEnv function using
// production environment.
func ValidateThirdParty(initData string, botID int64, expIn time.Duration) error {
	return ValidateThirdPartyWithEnv(initData, botID, expIn, false)
}
