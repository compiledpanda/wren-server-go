package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestE2EMethodNotAllowed(t *testing.T) {
	ts := httptest.NewServer(Routes())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/bogus")
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	exp := `{"code":"METHOD_NOT_ALLOWED","description":"Method Not Allowed"}
`

	if string(body) != exp {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}
