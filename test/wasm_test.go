//+build js,wasm

package wasm_test

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFetch(t *testing.T) {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Unexpected StatusCode %d", resp.StatusCode)
	}

	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(rb) != "stuff" {
		t.Errorf("Unexpected Body: %q", string(rb))
	}
}
