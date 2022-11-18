package initdata

import (
	"reflect"
	"testing"
	"time"
)

const (
	defaultInitData = "query_id=AAHdF6IQAAAAAN0XohDhrOrc&user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%7D&auth_date=1662771648&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2"
	token           = "5768337691:AAH5YkoiEuPk8-FZa32hStHTqXiLPtAEhx8"
)

type testValidate struct {
	initData    string
	expIn       time.Duration
	expectedErr error
}

type testParse struct {
	initData    string
	expectedErr error
	expectedRes *InitData
}

var testsValidate = []testValidate{
	{
		initData: defaultInitData,
	},
	{
		initData:    defaultInitData,
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
		initData:    defaultInitData + "abc",
		expectedErr: ErrSignInvalid,
	},
}

var testsParse = []testParse{
	{
		initData:    defaultInitData + ";",
		expectedErr: ErrUnexpectedFormat,
	},
	{
		initData: defaultInitData,
		expectedRes: &InitData{
			QueryID: "AAHdF6IQAAAAAN0XohDhrOrc",
			User: &User{
				Id:           279058397,
				FirstName:    "Vladislav",
				LastName:     "Kibenko",
				Username:     "vdkfrost",
				LanguageCode: "ru",
				IsPremium:    true,
			},
			CanSendAfterRaw: 0,
			AuthDateRaw:     1662771648,
			Hash:            "c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2",
		},
	},
}

func TestValidate(t *testing.T) {
	for _, test := range testsValidate {
		if err := Validate(test.initData, token, test.expIn); err != test.expectedErr {
			t.Errorf("expected error to be %q. Received %q", test.expectedErr, err)
		}
	}
}

func TestParse(t *testing.T) {
	for _, test := range testsParse {
		if data, err := Parse(test.initData); err != test.expectedErr {
			t.Errorf("expected error to be %q. Received %q", test.expectedErr, err)
		} else if !reflect.DeepEqual(data, test.expectedRes) {
			t.Errorf("expected result to be %+v. Received %+v", test.expectedRes, data)
		}
	}
}
