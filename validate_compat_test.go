package initdata

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// validateJoinedPairs is the previous Validate implementation. It is kept to
// check that Validate returns the same result.
func validateJoinedPairs(initData, token string, expIn time.Duration) error {
	q, err := parseValidationQuery(initData)
	if err != nil {
		return err
	}

	var (
		authDate time.Time
		hash     string
		pairs    = make([]string, 0, len(q))
	)

	for k, v := range q {
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
		pairs = append(pairs, k+"="+v[0])
	}

	if hash == "" {
		return ErrSignMissing
	}

	if expIn > 0 {
		if authDate.IsZero() {
			return ErrAuthDateMissing
		}

		if authDate.Add(expIn).Before(time.Now()) {
			return ErrExpired
		}
	}

	sort.Strings(pairs)

	if signJoined(strings.Join(pairs, "\n"), token) != hash {
		return ErrSignInvalid
	}
	return nil
}

var _validateErrors = []error{
	ErrUnexpectedFormat,
	ErrSignMissing,
	ErrAuthDateMissing,
	ErrAuthDateInvalid,
	ErrExpired,
	ErrSignInvalid,
	strconv.ErrSyntax,
	strconv.ErrRange,
}

func checkValidateCompat(t *testing.T, initData string, expIn time.Duration) {
	t.Helper()

	want := validateJoinedPairs(initData, validateTestToken, expIn)
	got := Validate(initData, validateTestToken, expIn)

	if (want == nil) != (got == nil) {
		t.Fatalf("Validate(%q) = %v, want %v", initData, got, want)
	}
	for _, target := range _validateErrors {
		if errors.Is(want, target) != errors.Is(got, target) {
			t.Fatalf("Validate(%q) = %v, want %v", initData, got, want)
		}
	}
}

func signedTestInitData(payload map[string]string, authDate time.Time) string {
	q := url.Values{}
	for k, v := range payload {
		q.Set(k, v)
	}
	q.Set("auth_date", strconv.FormatInt(authDate.Unix(), 10))
	q.Set("hash", Sign(payload, validateTestToken, authDate))
	return q.Encode()
}

func validateCompatSeeds() []string {
	now := time.Now()
	return []string{
		validateTestInitData,
		validateTestInitData + "abc",
		validateTestInitData + ";",
		validateTestInitData + "&query_id=again",
		"",
		"no_hash=true",
		"hash=abc",
		"hash=abc&auth_date=test",
		"hash=abc&auth_date=99999999999999999999",
		"hash=abc&auth_date=1662771917",
		"%zz=1&hash=abc",
		signedTestInitData(map[string]string{"user": `{"id":1}`}, now),
		signedTestInitData(map[string]string{"user": `{"id":1}`}, now.Add(-time.Hour)),
		signedTestInitData(map[string]string{"a": "1", "a1": "2", "a-b": "3", "a_b": "4", "a=": "5"}, now),
		signedTestInitData(map[string]string{"user": "a+b c&d=e%", "start_param": "x\ny"}, now),
		strings.ToUpper(signedTestInitData(map[string]string{"user": `{"id":1}`}, now)),
	}
}

func TestValidateCompat(t *testing.T) {
	for _, initData := range validateCompatSeeds() {
		for _, expIn := range []time.Duration{0, time.Minute, 24 * time.Hour} {
			checkValidateCompat(t, initData, expIn)
		}
	}
}

func FuzzValidateCompat(f *testing.F) {
	for _, initData := range validateCompatSeeds() {
		f.Add(initData)
	}
	f.Fuzz(func(t *testing.T, initData string) {
		checkValidateCompat(t, initData, 0)
		checkValidateCompat(t, initData, time.Hour)
	})
}

func FuzzValidateSignedCompat(f *testing.F) {
	f.Add(`{"id":1}`, "abc", "query_id", "AAH")
	f.Add("", "", "a1", "x")
	f.Fuzz(func(t *testing.T, user, startParam, key, value string) {
		payload := map[string]string{"user": user, "start_param": startParam}
		if key != "hash" && key != "auth_date" {
			payload[key] = value
		}
		initData := signedTestInitData(payload, time.Now())
		if err := Validate(initData, validateTestToken, time.Hour); err != nil {
			if validateJoinedPairs(initData, validateTestToken, time.Hour) == nil {
				t.Fatalf("Validate(%q) = %v, want nil", initData, err)
			}
		}
		checkValidateCompat(t, initData, time.Hour)
	})
}

func TestValidateDifferentTokens(t *testing.T) {
	const otherToken = "1234567890:AAbbccddeeffgghhiijjkkllmmnnooppqq"

	authDate := time.Now()
	payload := map[string]string{"user": `{"id":1}`}
	own := signedTestInitData(payload, authDate)

	q := url.Values{}
	q.Set("user", `{"id":1}`)
	q.Set("auth_date", strconv.FormatInt(authDate.Unix(), 10))
	q.Set("hash", Sign(payload, otherToken, authDate))
	other := q.Encode()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				if err := Validate(own, validateTestToken, 0); err != nil {
					t.Errorf("own token: %v", err)
				}
				if err := Validate(other, otherToken, 0); err != nil {
					t.Errorf("other token: %v", err)
				}
				if err := Validate(other, validateTestToken, 0); !errors.Is(err, ErrSignInvalid) {
					t.Errorf("wrong token: %v", err)
				}
			}
		}()
	}
	wg.Wait()
}

func BenchmarkValidate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := Validate(validateTestInitData, validateTestToken, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateJoinedPairs(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := validateJoinedPairs(validateTestInitData, validateTestToken, 0); err != nil {
			b.Fatal(err)
		}
	}
}

// signJoined is the previous sign implementation.
func signJoined(payload, key string) string {
	skHmac := hmac.New(sha256.New, []byte("WebAppData"))
	skHmac.Write([]byte(key))

	impHmac := hmac.New(sha256.New, skHmac.Sum(nil))
	impHmac.Write([]byte(payload))

	return hex.EncodeToString(impHmac.Sum(nil))
}
