package test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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

func VerifyStringResponse(t *testing.T, rr *httptest.ResponseRecorder, statusCode int, body string) {
	if status := rr.Code; status != statusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, statusCode)
	}

	if rr.Body.String() != body {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), body)
	}
}
