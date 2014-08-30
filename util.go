package douban

import (
	"strings"
)

func Urlencode(base_url string, params map[string]string) string {
	_b := "?"
	for k, v := range params {
		_b += strings.Join([]string{k, "=", v, "&"}, "")
	}
	return base_url + _b[:len(_b)-1]
}
