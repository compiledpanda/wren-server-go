package server

import "testing"

func TestSetup(t *testing.T) {
	srv, _ := Setup(&Config{}) // TODO test error

	// TODO #2 Test Configuration

	exp := "0.0.0.0:8985"

	if srv.Addr != exp {
		t.Errorf("Unexpected Addr returned: got %v want %v",
			srv.Addr, exp)
	}

}
