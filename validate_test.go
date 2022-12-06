package initdata

import (
	"testing"
	"time"
)

const (
	validateTestInitData = "query_id=AAHdF6IQAAAAAN0XohDhrOrc&user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%7D&auth_date=1662771648&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2"
	validateTestToken    = "5768337691:AAH5YkoiEuPk8-FZa32hStHTqXiLPtAEhx8"
)

type testValidate struct {
	initData    string
	expIn       time.Duration
	expectedErr error
}

var testsValidate = []testValidate{
	{
		initData: validateTestInitData,
	},
	{
		initData:    validateTestInitData,
		expectedErr: ErrExpired,
		expIn:       time.Second,
	},
	{
		initData:    "here comes something wrong;",
		expectedErr: ErrUnexpectedFormat,
	},
	{
		initData:    "no_hash=true",
		expectedErr: ErrSignMissing,
	},
	{
		initData:    "hash=abc",
		expectedErr: ErrAuthDateMissing,
		expIn:       time.Second,
	},
	{
		initData:    "hash=abc&auth_date=1662771917",
		expectedErr: ErrExpired,
		expIn:       time.Second,
	},
	{
		initData:    validateTestInitData + "abc",
		expectedErr: ErrSignInvalid,
	},
}

func TestValidate(t *testing.T) {
	for _, test := range testsValidate {
		if err := Validate(test.initData, validateTestToken, test.expIn); err != test.expectedErr {
			t.Errorf("expected error to be %q. Received %q", test.expectedErr, err)
		}
	}
}
