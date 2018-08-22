package pfclient

import (
	"testing"
)

func TestNewContainerListFromByte(t *testing.T) {
	tables := []struct {
		hostname string
		image    string
	}{
		{"test-01", "16.04"},
		{"test-02", "16.04"},
		{"test-03", "16.04"},
	}

	b := []byte(`{
		"api_version": "1.0",
		"data": {
			"items": [
				{"hostname": "test-01", "image": "16.04"},
				{"hostname": "test-02", "image": "16.04"},
				{"hostname": "test-03", "image": "16.04"}
			]
		}
	}`)
	cl, _ := NewContainerListFromByte(b)

	for i, table := range tables {
		if cl.Get(i).Hostname != table.hostname {
			t.Errorf("Incorrect container hostname generated, got: %s, want: %s.",
				cl.Get(i).Hostname,
				table.hostname)
		}

		if cl.Get(i).Image != table.image {
			t.Errorf("Incorrect container image generated, got: %s, want: %s.",
				cl.Get(i).Image,
				table.image)
		}
	}
}
