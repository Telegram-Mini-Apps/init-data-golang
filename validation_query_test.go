package initdata

import (
	"errors"
	"testing"
)

func TestParseValidationQueryRejectsDuplicateDecodedKeys(t *testing.T) {
	tests := []struct {
		name      string
		initData  string
		wantError bool
	}{
		{
			name:     "unique keys",
			initData: "auth_date=1700000000&hash=abc",
		},
		{
			name:      "duplicate field",
			initData:  "query_id=first&query_id=second",
			wantError: true,
		},
		{
			name:      "duplicate after key decoding",
			initData:  "query_id=first&query%5Fid=second",
			wantError: true,
		},
		{
			name:      "duplicate hash",
			initData:  "hash=first&hash=second",
			wantError: true,
		},
		{
			name:      "duplicate signature",
			initData:  "signature=first&signature=second",
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseValidationQuery(test.initData)
			if test.wantError && !errors.Is(err, ErrUnexpectedFormat) {
				t.Fatalf("expected %q, got %v", ErrUnexpectedFormat, err)
			}
			if !test.wantError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestValidatorsRejectDuplicateDecodedKeys(t *testing.T) {
	validators := []struct {
		name     string
		data     string
		validate func(string) error
	}{
		{"HMAC", validateTestInitData, func(data string) error { return Validate(data, validateTestToken, 0) }},
		{"Ed25519", _validateThirdPartyTestInitData, func(data string) error { return ValidateThirdParty(data, _validateThirdPartyBotID, 0) }},
	}
	for _, validator := range validators {
		t.Run(validator.name, func(t *testing.T) {
			if err := validator.validate(validator.data); err != nil {
				t.Fatalf("valid baseline rejected: %v", err)
			}
			for _, suffix := range []string{"&user=other", "&%75ser=other", "&hash=other"} {
				t.Run(suffix, func(t *testing.T) {
					if err := validator.validate(validator.data + suffix); !errors.Is(err, ErrUnexpectedFormat) {
						t.Fatalf("expected ErrUnexpectedFormat for duplicate decoded key, got %v", err)
					}
				})
			}
		})
	}
	if err := ValidateThirdParty(_validateThirdPartyTestInitData+"&signature=other", _validateThirdPartyBotID, 0); !errors.Is(err, ErrUnexpectedFormat) {
		t.Fatalf("expected ErrUnexpectedFormat for duplicate signature, got %v", err)
	}
}
