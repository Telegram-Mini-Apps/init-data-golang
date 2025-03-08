package initdata

import (
	"errors"
	"testing"
	"time"
)

const (
	signTestInitData = "query_id=AAHdF6IQAAAAAN0XohDhrOrc&user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%7D&auth_date=1662771648"
	signTestToken    = "5768337691:AAH5YkoiEuPk8-FZa32hStHTqXiLPtAEhx8"
	signTestHash     = "c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2"
)

var (
	signTestAuthDate = time.Unix(1662771648, 0)
)

type testSignQueryString struct {
	initData    string
	expectedErr error
}

var testsSignQueryString = []testSignQueryString{
	// Default stable case. No fields to sanitize.
	{
		initData: "query_id=AAHdF6IQAAAAAN0XohDhrOrc&user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%7D",
	},

	// Should ignore "auth_date" and "hash" fields.
	{
		initData: "query_id=AAHdF6IQAAAAAN0XohDhrOrc&user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%7D&auth_date=1662771648&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2",
	},
}

func TestSignQueryString(t *testing.T) {
	for _, test := range testsSignQueryString {
		hash, err := SignQueryString(test.initData, signTestToken, signTestAuthDate)
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("expected error to be %q. Received %q", test.expectedErr, err)
		}
		if hash != signTestHash {
			t.Errorf("expected result to be %+v. Received %+v", signTestHash, hash)
		}
	}
}
