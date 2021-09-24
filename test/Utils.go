package test

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

/** Test Utilities **/
func CallHandler(t *testing.T, handler func(http.ResponseWriter, *http.Request), method string, url string, body io.Reader) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler)

	h.ServeHTTP(rr, req)

	return rr
}

func VerifyStringResponse(t *testing.T, expectedStatusCode int, expectedBody string, statusCode int, body string) {
	if expectedStatusCode != statusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			expectedStatusCode, statusCode)
	}

	if expectedBody != body {
		t.Errorf("handler returned unexpected body: got %v want %v",
			expectedBody, body)
	}
}

func CallGetEndpoint(t *testing.T, router *mux.Router, url string) (statusCode int, headers http.Header, body []byte) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + url)
	if err != nil {
		t.Fatal(err)
	}
	statusCode = res.StatusCode
	headers = res.Header

	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	return
}

func VerifyHeader(t *testing.T, header string, expectedHeader string) {
	if header != expectedHeader {
		t.Errorf("Header does not match: got %v want %v",
			header, expectedHeader)
	}
}

func VerifyDigestHeader(t *testing.T, header string, expectedContents []byte) {
	// Calculate the hash
	hasher := sha256.New()
	_, err := hasher.Write(expectedContents)
	if err != nil {
		t.Fatal(err)
	}
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	expectedHeader := "sha-256=" + hash

	if header != expectedHeader {
		t.Errorf("Digest does not match: got %v want %v",
			header, expectedHeader)
	}
}
