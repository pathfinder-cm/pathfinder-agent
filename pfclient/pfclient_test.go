package pfclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchContainersFromServer(t *testing.T) {
	node := "test-01"
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

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write(b)
	}))
	defer func() { testServer.Close() }()

	pfclient := NewPfclient(&http.Client{}, testServer.URL, "", "")
	cl, _ := pfclient.FetchContainersFromServer(node)
	for i, table := range tables {
		if (*cl)[i].Hostname != table.hostname {
			t.Errorf("Incorrect container hostname fetched, got: %s, want: %s.",
				(*cl)[i].Hostname,
				table.hostname)
		}

		if (*cl)[i].Image != table.image {
			t.Errorf("Incorrect container image fetched, got: %s, want: %s.",
				(*cl)[i].Image,
				table.image)
		}
	}
}

func TestMarkContainerAsProvisioned(t *testing.T) {
	tables := []struct {
		node     string
		hostname string
	}{
		{"test-01", "test-c-01"},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))
	defer func() { testServer.Close() }()

	pfclient := NewPfclient(&http.Client{}, testServer.URL, "", "")
	ok, _ := pfclient.MarkContainerAsProvisioned(tables[0].node, tables[0].hostname)
	if ok != true {
		t.Errorf("Error when marking container as provisioned")
	}
}
