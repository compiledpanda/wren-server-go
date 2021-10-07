package test

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

/** Test Utilities **/
func CallHandler(t *testing.T, handler func(http.ResponseWriter, *http.Request), method string, url string, headers map[string]string, body []byte) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(handler)

	h.ServeHTTP(rr, req)

	return rr
}

func CallGetEndpoint(t *testing.T, router *mux.Router, url string) *http.Response {
	ts := httptest.NewServer(router)
	defer ts.Close()

	res, err := http.Get(ts.URL + url)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func VerifyStringResponse(t *testing.T, res *http.Response, status int, body string) {
	if res.StatusCode != status {
		t.Errorf("Status code does not match: got %v want %v",
			res.StatusCode, status)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "text/plain" {
		t.Error("Content Type is not text/plain")
	}

	actualBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}
	if string(actualBody) != body {
		t.Errorf("Body does not match: got %v want %v", string(actualBody), body)
	}
}

func VerifyJSONResponse(t *testing.T, res *http.Response, status int, body string) {
	if res.StatusCode != status {
		t.Errorf("Status code does not match: got %v want %v",
			res.StatusCode, status)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Error("Content Type is not application/json")
	}

	actualBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}
	if string(actualBody) != body {
		t.Errorf("Body does not match: got %v want %v", string(actualBody), body)
	}
}

func VerifyByteResponse(t *testing.T, res *http.Response, status int, body []byte) {
	if res.StatusCode != status {
		t.Errorf("Status code does not match: got %v want %v",
			res.StatusCode, status)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/octet-stream" {
		t.Error("Content Type is not application/octet-stream")
	}

	digest := res.Header.Get("Digest")
	hasher := sha256.New()
	_, err := hasher.Write(body)
	if err != nil {
		t.Fatal(err)
	}
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	expectedDigest := "sha-256=" + hash

	if digest != expectedDigest {
		t.Errorf("Digest does not match: got %v want %v", digest, expectedDigest)
	}

	actualBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actualBody, body) {
		t.Errorf("Body does not match: got %v want %v", actualBody, body)
	}
}

func VerifyRecordedStringResponse(t *testing.T, rr *httptest.ResponseRecorder, status int, body string) {
	if rr.Code != status {
		t.Errorf("Status code does not match: got %v want %v",
			rr.Code, status)
	}

	contentType := rr.Result().Header.Get("Content-Type")
	if contentType != "text/plain" {
		t.Error("Content Type is not text/plain")
	}

	actualBody := rr.Body.String()
	if actualBody != body {
		t.Errorf("Body does not match: got %v want %v", actualBody, body)
	}
}

func VerifyRecordedJSONResponse(t *testing.T, rr *httptest.ResponseRecorder, status int, body string) {
	if rr.Code != status {
		t.Errorf("Status code does not match: got %v want %v",
			rr.Code, status)
	}

	contentType := rr.Result().Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Error("Content Type is not application/json")
	}

	actualBody := rr.Body.String()
	if actualBody != body {
		t.Errorf("Body does not match: got %v want %v", actualBody, body)
	}
}

func VerifyRecordedByteResponse(t *testing.T, rr *httptest.ResponseRecorder, status int, body []byte) {
	if rr.Code != status {
		t.Errorf("Status code does not match: got %v want %v",
			rr.Code, status)
	}

	contentType := rr.Result().Header.Get("Content-Type")
	if contentType != "application/octet-stream" {
		t.Error("Content Type is not application/octet-stream")
	}

	digest := rr.Result().Header.Get("Digest")
	hasher := sha256.New()
	_, err := hasher.Write(body)
	if err != nil {
		t.Fatal(err)
	}
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	expectedDigest := "sha-256=" + hash

	if digest != expectedDigest {
		t.Errorf("Digest does not match: got %v want %v", digest, expectedDigest)
	}

	actualBody := rr.Body.Bytes()
	if !bytes.Equal(actualBody, body) {
		t.Errorf("Body does not match: got %v want %v", actualBody, body)
	}
}

func VerifyErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, status int, code string) {
	type serverError struct {
		Code        string `json:"code"`
		Description string `json:"description"`
	}

	if rr.Code != status {
		t.Errorf("Status code does not match: got %v want %v",
			rr.Code, status)
	}

	contentType := rr.Result().Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Error("Content Type is not application/json")
	}

	var e serverError
	err := json.Unmarshal(rr.Body.Bytes(), &e)

	if err != nil {
		t.Errorf("Unable to unmarshal result: %v", err)
	}

	if e.Code != code {
		t.Errorf("Error code does not match: got %v want %v",
			e.Code, code)
	}

	if e.Description == "" {
		t.Error("Error description is empty")
	}
}
