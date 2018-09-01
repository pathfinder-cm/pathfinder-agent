package pfclient

import (
	"testing"
)

func TestNewRegisterFromByte(t *testing.T) {
	tables := []struct {
		hostname            string
		authenticationToken string
	}{
		{"test-01", "abc"},
	}

	b := []byte(`{
		"api_version": "1.0",
		"data": {
			"hostname": "test-01",
			"authentication_token": "abc"
		}
	}`)
	register, _ := NewRegisterFromByte(b)

	for _, table := range tables {
		if register.Hostname != table.hostname {
			t.Errorf("Incorrect hostname generated, got: %s, want: %s.",
				register.Hostname,
				table.hostname)
		}

		if register.AuthenticationToken != table.authenticationToken {
			t.Errorf("Incorrect authentication token generated, got: %s, want: %s.",
				register.AuthenticationToken,
				table.authenticationToken)
		}
	}
}
