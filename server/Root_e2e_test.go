package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestE2ERoot(t *testing.T) {
	ts := httptest.NewServer(Routes())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	exp := "Hello World\n"

	if string(body) != exp {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}
