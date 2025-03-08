package initdata

import (
	"errors"
	"testing"
	"time"
)

const (
	_validateThirdPartyTestInitData = "user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%20%2B%20-%20%3F%20%5C%2F%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg%22%7D&chat_instance=8134722200314281151&chat_type=private&auth_date=1733584787&hash=2174df5b000556d044f3f020384e879c8efcab55ddea2ced4eb752e93e7080d6&signature=zL-ucjNyREiHDE8aihFwpfR9aggP2xiAo3NSpfe-p7IbCisNlDKlo7Kb6G4D0Ao2mBrSgEk4maLSdv6MLIlADQ"
	_validateThirdPartyBotID        = 7342037359
)

func TestValidateThirdPartyValid(t *testing.T) {
	err := ValidateThirdParty(_validateThirdPartyTestInitData, _validateThirdPartyBotID, 0)
	if err != nil {
		t.Errorf("expected not to return error. Received: %q", err)
	}
}

func TestValidateThirdPartyExpired(t *testing.T) {
	err := ValidateThirdParty(_validateThirdPartyTestInitData, _validateThirdPartyBotID, time.Second)
	if !errors.Is(err, ErrExpired) {
		t.Errorf("expected to receive %q. Received: %q", ErrExpired, err)
	}
}

func TestValidateThirdPartyInvalidFormat(t *testing.T) {
	err := ValidateThirdParty("here comes something wrong;", 1, 0)
	if !errors.Is(err, ErrUnexpectedFormat) {
		t.Errorf("expected to receive %q. Received: %q", ErrUnexpectedFormat, err)
	}
}

func TestValidateThirdPartySignMissing(t *testing.T) {
	err := ValidateThirdParty("no_signature=true", 1, 0)
	if !errors.Is(err, ErrSignMissing) {
		t.Errorf("expected to receive %q. Received: %q", ErrSignMissing, err)
	}
}

func TestValidateThirdPartyAuthDateMissing(t *testing.T) {
	err := ValidateThirdParty("signature=abc", 1, time.Second)
	if !errors.Is(err, ErrAuthDateMissing) {
		t.Errorf("expected to receive %q. Received: %q", ErrAuthDateMissing, err)
	}
}

func TestValidateThirdPartyAuthDateInvalid(t *testing.T) {
	err := ValidateThirdParty("signature=abc&auth_date=test", 1, time.Second)
	if !errors.Is(err, ErrAuthDateInvalid) {
		t.Errorf("expected to receive %q. Received: %q", ErrAuthDateInvalid, err)
	}
}

func TestValidateThirdPartySignInvalid(t *testing.T) {
	err := ValidateThirdParty(_validateThirdPartyTestInitData+"a", 1, 0)
	if !errors.Is(err, ErrSignInvalid) {
		t.Errorf("expected to receive %q. Received: %q", ErrSignInvalid, err)
	}
}
