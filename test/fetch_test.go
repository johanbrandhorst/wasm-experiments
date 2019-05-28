package fetch_test

import (
	"encoding/json"
	"net/http"
	"testing"
)

type jsonResp struct {
	Headers struct {
		Accept         string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		Host           string `json:"Host"`
		Origin         string `json:"Origin"`
		Referer        string `json:"Referer"`
		UserAgent      string `json:"User-Agent"`
	} `json:"headers"`
	Origin string `json:"origin"`
	URL    string `json:"url"`
}

func TestFetch(t *testing.T) {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected StatusCode %d", resp.StatusCode)
	}

	var r jsonResp
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		t.Fatal(err)
	}

	if r.URL != "https://httpbin.org/get" {
		t.Errorf("Unexpected request URL: %q", r.URL)
	}
}
