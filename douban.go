package douban

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	return &Client{API_KEY, API_SECRET, redrict_url, scope, nil}
}

func (c *Client) get(url string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", API_HOST+url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.access_token["access_token"].(string))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return content
}

func (c *Client) AuthorizeUrl() string {
	return Urlencode(AUTH_HOST+"/service/auth2/auth", map[string]string{
		"client_id":     c.APT_KEY,
		"redirect_uri":  c.redirct_url,
		"response_type": "code",
		"scope":         strings.Join(c.scope, ","),
	})
}

func (c *Client) AuthWithToken(token map[string]interface{}) {
	c.access_token = token
}

func (c *Client) AuthWithCode(code string) {
	resp, err := http.PostForm(AUTH_HOST+"/service/auth2/token", url.Values{
		"client_id":     {c.APT_KEY},
		"client_secret": {c.API_SECRET},
		"redirect_uri":  {c.redirct_url},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	})
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(content))
	json.Unmarshal(content, &c.access_token)
}

func (c *Client) RefreshToken() {
	resp, err := http.PostForm(AUTH_HOST+"/service/auth2/token", url.Values{
		"client_id":     {c.APT_KEY},
		"client_secret": {c.API_SECRET},
		"redirect_uri":  {c.redirct_url},
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.access_token["refresh_token"].(string)},
	})
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(content, &c.access_token)
}

// User

func (c *Client) Me() []byte {
	return c.get("/v2/user/~me")
}

func (c *Client) GetUserbyId(id int32) []byte {
	return c.get("/v2/user/" + string(id))
}

func (c *Client) GetUserbyName(username string) []byte {
	return c.get("/v2/user/" + username)
}

func (c *Client) SearchUserByKeywords(keywords []string, start, count int32) []byte {
	return c.get(Urlencode("/v2/user", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": string(start),
		"count": string(count),
	}))
}

// Book
func (c *Client) GetBookById(id int32) []byte {
	return c.get("/v2/book/" + string(id))
}

func (c *Client) GetBookByISBN(isbn string) []byte {
	return c.get("/v2/book/isbn/" + isbn)
}

func (c *Client) SearchBookByKeywords(keywords []string, start, count int32) []byte {
	return c.get(Urlencode("/v2/book/search", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": string(start),
		"count": string(count),
	}))
}

func (c *Client) SearchBookByTag(tag string, start, count int32) []byte {
	return c.get(Urlencode("/v2/book/search", map[string]string{
		"tag":   tag,
		"start": string(start),
		"count": string(count),
	}))
}

func (c *Client) GetTagsOfBookById(id int32) []byte {
	return c.get("/v2/book/" + string(id) + "/tags")
}

/*
http://developers.douban.com/wiki/?title=book_v2#get_user_book_tags
func (c *Client) GetTagsOfBooksOfUser() {

}
*/

// http://developers.douban.com/wiki/?title=book_v2#get_user_collections
func (c *Client) GetBookCollectionsOfUser(username string) []byte {
	return c.get("/v2/book/user/" + username + "/collections")
}

// http://developers.douban.com/wiki/?title=book_v2#get_book_collection
func (c *Client) GetBookCollectionOfUser(book_id, username string) []byte {
	return c.get(Urlencode("/v2/book/"+book_id+"/collection", map[string]string{
		"user_id": username,
	}))
}

// http://developers.douban.com/wiki/?title=book_v2#get_user_annotations
func (c *Client) GetUserAnnotations(username string) {
	return c.get("/v2/book/user/" + username + "/annotations")
}

// http://developers.douban.com/wiki/?title=book_v2#get_book_annotations
func (c *Client) GetBookAnnotations(book_id int32) {
	return c.get("/v2/book/" + string(book_id) + "/annotations")
}

// http://developers.douban.com/wiki/?title=book_v2#get_annotation
func (c *Client) GetAnnotationById(id int32) {
	return c.get("/v2/book/annotation/" + string(id))
}
