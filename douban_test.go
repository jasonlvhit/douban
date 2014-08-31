package douban

import (
	"encoding/json"
	"testing"
)

const (
	TEST_API_KEY      = "08bceb66ac35161b1a36c234049f4251"
	TEST_API_SECRET   = "4ec658d9e4549ce7"
	TEST_REDIRECT_URI = "http://feedscard.sinaapp.com"
)

var scopes = map[string]([]string){
	//豆瓣公共
	"common": {"douban_basic_common"},
	//东西
	"dongxi":    {"commodity_basic_r", "commodity_basic_w"},
	"movie":     {"movie_basic", "movie_basic_r", "movie_basic_w"},
	"travel":    {"travel_basic_r"},
	"community": {"community_basic_note", "community_basic_user", "community_basic_photo", "community_basic_online"},
	"thing":     {"thing_basic_r", "thing_basic_w"},
	"book":      {"book_basic_r", "book_basic_w"},
	"music":     {"music_basic_r", "music_basic_w", "music_artist_r"},
	"shuo":      {"shuo_basic_r", "shuo_basic_w"},
	//同城
	"event": {"event_basic_r", "event_basic_w", "event_drama_r", "event_drama_w"},
}

var client = NewClient(TEST_API_KEY, TEST_API_SECRET, TEST_REDIRECT_URI, scopes["common"])

func TestNewClient(t *testing.T) {
	if client.API_SECRET != TEST_API_KEY {
		t.Errorf("API_KEY want %s, got %s", TEST_API_KEY, client.API_KEY)
	}
	if client.API_SECRET != TEST_API_SECRET {
		t.Errorf("")
	}
}

func TestAuthorizeUrl(t *testing.T) {
	if client.AuthorizeUrl() != "https://www.douban.com/service/auth2/auth?redirect_uri=http://feedscard.sinaapp.com&scope=douban_basic_common&response_type=code&client_id=08bceb66ac35161b1a36c234049f4251" {
		t.Errorf()
	}
}

type genericJson map[string]interface{}

func TestMe(t *testing.T) {
	me := client.Me()
	var v genericJson
	json.Unmarshal(me, &v)
	if v["name"] != "辣椒面儿" {
		t.Errorf("Client Me() get %s, want %s", v["name"], "辣椒面儿")
	}
}
