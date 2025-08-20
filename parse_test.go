package initdata

import (
	"errors"
	"reflect"
	"testing"
)

const (
	_parseTestInitData = "query_id=AAHdF6IQAAAAAN0XohDhrOrc&user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%7D&auth_date=1662771648&hash=c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2&start_param=abc"
)

type testParse struct {
	initData    string
	expectedErr error
	expectedRes *InitData
}

var testsParse = []testParse{
	{
		initData:    _parseTestInitData + ";",
		expectedErr: ErrUnexpectedFormat,
	},
	{
		initData: _parseTestInitData,
		expectedRes: &InitData{
			QueryID: "AAHdF6IQAAAAAN0XohDhrOrc",
			User: &User{
				ID:           279058397,
				FirstName:    "Vladislav",
				LastName:     "Kibenko",
				Username:     "vdkfrost",
				LanguageCode: "ru",
				IsPremium:    true,
			},
			CanSendAfterRaw: 0,
			AuthDateRaw:     1662771648,
			StartParam:      "abc",
			Hash:            "c501b71e775f74ce10e377dea85a7ea24ecd640b223ea86dfe453e0eaed2e2b2",
		},
	},
}

func TestParse(t *testing.T) {
	for _, test := range testsParse {
		if data, err := Parse(test.initData); !errors.Is(err, test.expectedErr) {
			t.Errorf("expected error to be %q. \nReceived %q", test.expectedErr, err)
		} else if !reflect.DeepEqual(data, test.expectedRes) {
			t.Errorf("expected result to be %+v. \nReceived %+v", test.expectedRes, data)
		}
	}
}
