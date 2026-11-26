package initdata

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// parseJSONRebuild is the previous Parse implementation. It is kept to check
// that Parse decodes init data the same way.
func parseJSONRebuild(initData string) (*InitData, error) {
	q, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("parse init data as query: %w: %w", err, ErrUnexpectedFormat)
	}

	pairs := make([]string, 0, len(q))
	for k, v := range q {
		val := v[0]
		valFormat := "%q:%q"
		if isString := _stringProps[k]; !isString && json.Valid([]byte(val)) {
			valFormat = "%q:%s"
		}
		pairs = append(pairs, fmt.Sprintf(valFormat, k, val))
	}

	var d InitData
	jStr := fmt.Sprintf("{%s}", strings.Join(pairs, ","))
	if err := json.Unmarshal([]byte(jStr), &d); err != nil {
		return nil, fmt.Errorf("unmarshal init data: %w: %w", err, ErrUnexpectedFormat)
	}
	return &d, nil
}

var _parseCompatSeeds = []string{
	_parseTestInitData,
	_parseTestInitData + ";",
	_parseTestInitData + "&user=%7B%22id%22%3A1%7D",
	"",
	"%zz=1",
	"auth_date=abc",
	"auth_date=%201",
	"auth_date=1.5",
	"auth_date=%221%22",
	"can_send_after=10&auth_date=1",
	"chat_instance=-1234567890123456789",
	"chat_instance=abc",
	"chat_type=sender",
	"chat_type=%22sender%22",
	"chat_type=123",
	"chat_type=null",
	"hash=00ff",
	"hash=123",
	"query_id=AAHdF6IQAAAAAN0XohDhrOrc",
	"start_param=123",
	"start_param=%7B%22a%22%3A1%7D",
	"start_param=",
	"user=null",
	"user=not%20json",
	"user=%7B%22ID%22%3A5%2C%22First_Name%22%3A%22a%22%7D",
	"Auth_Date=5&HASH=x",
	"receiver=%7B%22id%22%3A2%7D&chat=%7B%22id%22%3A3%2C%22type%22%3A%22group%22%2C%22title%22%3A%22t%22%7D",
	"signature=abc&unknown=1",
}

func hasFoldedDuplicateKeys(initData string) bool {
	q, err := url.ParseQuery(initData)
	if err != nil {
		return false
	}

	seen := make(map[string]bool, len(q))
	for k := range q {
		lower := strings.ToLower(k)
		if seen[lower] {
			return true
		}
		seen[lower] = true
	}
	return false
}

func checkParseCompat(t *testing.T, initData string) {
	t.Helper()

	// The previous implementation depended on map iteration order when the
	// same property was passed with different letter cases.
	if hasFoldedDuplicateKeys(initData) {
		return
	}

	want, wantErr := parseJSONRebuild(initData)
	got, gotErr := Parse(initData)

	if wantErr == nil && gotErr != nil {
		t.Fatalf("Parse(%q) failed: %v", initData, gotErr)
	}
	if wantErr == nil && !reflect.DeepEqual(want, got) {
		t.Fatalf("Parse(%q) = %+v, want %+v", initData, got, want)
	}
}

func TestParseCompat(t *testing.T) {
	for _, initData := range _parseCompatSeeds {
		checkParseCompat(t, initData)
	}
}

func FuzzParseCompat(f *testing.F) {
	for _, initData := range _parseCompatSeeds {
		f.Add(initData)
	}
	f.Fuzz(checkParseCompat)
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Parse(_parseTestInitData); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseJSONRebuild(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := parseJSONRebuild(_parseTestInitData); err != nil {
			b.Fatal(err)
		}
	}
}
