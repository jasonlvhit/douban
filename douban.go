package douban

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const AUTH_HOST = "https://www.douban.com"
const API_HOST = "https://api.douban.com"

type Client struct {
	APT_KEY      string
	API_SECRET   string
	redirct_url  string
	scope        []string
	access_token map[string]interface{}
}

func NewClient(API_KEY, API_SECRET, redrict_url string, scope []string) *Client {
	return Client{API_KEY, API_SECRET, redrict_url, scope}
}

func (c *Client) AuthorizeUrl() string {
	return Urlencode(AUTH_HOST+"/service/auth2/auth", map[string]string{
		"client_id":     c.APT_KEY,
		"redirect_uri":  c.redirct_url,
		"response_type": "code",
		"scope":         strings.Join(c.scope, ","),
	})
}

func (c *Client) AuthWithToken() {
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", "https://api.douban.com/v2/user/~me", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.access_token["access_token"].(string))
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	// Do something with resq.
}

func (c *Client) AuthWithCode(code string) {
	resp, err := http.Post(Urlencode(AUTH_HOST+"service/auth2/token", map[string]string{
		"client_id":     c.APT_KEY,
		"client_secret": c.API_SECRET,
		"redirect_uri":  c.redirct_url,
		"grant_type":    "authorization_code",
		"code":          code,
	}), nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(c.access_token, &v)
	c.AuthWithToken()
	//resq
}

func (c *Client) RefreshToken() {
	resp, err := http.Post(Urlencode(AUTH_HOST+"service/auth2/token", map[string]string{
		"client_id":     c.APT_KEY,
		"client_secret": c.API_SECRET,
		"redirect_uri":  c.redirct_url,
		"grant_type":    "refresh_token",
		"refresh_token": c.access_token["refresh_token"],
	}), nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(c.access_token, &v)
}
