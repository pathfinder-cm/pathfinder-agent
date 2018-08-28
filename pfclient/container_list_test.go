package pfclient

import (
	"testing"
)

func TestNewContainerListFromByte(t *testing.T) {
	tables := []struct {
		hostname string
		image    string
		status   string
	}{
		{"test-01", "16.04", "SCHEDULED"},
		{"test-02", "16.04", "SCHEDULED"},
		{"test-03", "16.04", "SCHEDULED"},
	}

	b := []byte(`{
		"api_version": "1.0",
		"data": {
			"items": [
				{"hostname": "test-01", "image": "16.04", "status": "SCHEDULED"},
				{"hostname": "test-02", "image": "16.04", "status": "SCHEDULED"},
				{"hostname": "test-03", "image": "16.04", "status": "SCHEDULED"}
			]
		}
	}`)
	cl, _ := NewContainerListFromByte(b)

	for i, table := range tables {
		if (*cl)[i].Hostname != table.hostname {
			t.Errorf("Incorrect container hostname generated, got: %s, want: %s.",
				(*cl)[i].Hostname,
				table.hostname)
		}

		if (*cl)[i].Image != table.image {
			t.Errorf("Incorrect container image generated, got: %s, want: %s.",
				(*cl)[i].Image,
				table.image)
		}

		if (*cl)[i].Status != table.status {
			t.Errorf("Incorrect container status generated, got: %s, want: %s.",
				(*cl)[i].Status,
				table.status)
		}
	}
}
