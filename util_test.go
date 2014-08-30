package douban

import (
	"testing"
)

func TestUrlencode(t *testing.T) {
	base_url := "http://www.github.com"
	params := map[string]string{}
	url := Urlencode(base_url, params)
	if url != base_url {
		t.Errorf("Urlencode. empty params: %s, want %s", url, base_url)
	}
	params = map[string]string{
		"key1": "1",
		"key2": "2",
	}
	url = Urlencode(base_url, params)
	want := base_url + "?" + "key1=1&" + "key2=2"
	if url != want {
		t.Errorf("Urlencode. 2 pairs parmas: %s, want %s", url, want)
	}
}
