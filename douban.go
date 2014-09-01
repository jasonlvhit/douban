package douban

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	fmt.Println(API_HOST + url)
	req, err := http.NewRequest("GET", API_HOST+url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", c.access_token["access_token"].(string))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return content
}

/*

func (c *Client) post(purl string, data map[string]([]string)) []byte {
	http.Header.Set("jj", "ss")
	resp, err := http.PostForm(API_HOST+purl, data)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return content
}

*/

func (c *Client) post(purl string, data map[string]([]string)) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", API_HOST+purl, strings.NewReader(url.Values(data).Encode()))
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
		"scope":         strings.Join(c.scope, "+"),
	})
}

func (c *Client) AuthWithToken(token interface{}) {
	switch token.(type) {
	case string:
		c.access_token = map[string]interface{}{
			"access_token": token.(string),
		}
	case map[string]interface{}:
		c.access_token = token.(map[string]interface{})
	}
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

func (c *Client) GetUserbyId(id int) []byte {
	return c.get("/v2/user/" + strconv.Itoa(id))
}

func (c *Client) GetUserbyName(username string) []byte {
	return c.get("/v2/user/" + username)
}

func (c *Client) SearchUserByKeywords(keywords []string, start, count int) []byte {
	return c.get(Urlencode("/v2/user", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

// Book
func (c *Client) GetBookById(id int) []byte {
	return c.get("/v2/book/" + strconv.Itoa(id))
}

func (c *Client) GetBookByISBN(isbn string) []byte {
	return c.get("/v2/book/isbn/" + isbn)
}

func (c *Client) SearchBookByKeywords(keywords []string, start, count int) []byte {
	return c.get(Urlencode("/v2/book/search", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) SearchBookByTag(tag string, start, count int) []byte {
	return c.get(Urlencode("/v2/book/search", map[string]string{
		"tag":   tag,
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) GetTagsOfBookById(id int) []byte {
	return c.get("/v2/book/" + strconv.Itoa(id) + "/tags")
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
func (c *Client) GetUserAnnotations(username string) []byte {
	return c.get("/v2/book/user/" + username + "/annotations")
}

// http://developers.douban.com/wiki/?title=book_v2#get_book_annotations
func (c *Client) GetBookAnnotations(book_id int) []byte {
	return c.get("/v2/book/" + strconv.Itoa(book_id) + "/annotations")
}

// http://developers.douban.com/wiki/?title=book_v2#get_annotation
func (c *Client) GetAnnotationById(id int) []byte {
	return c.get("/v2/book/annotation/" + strconv.Itoa(id))
}

// Movie

func (c *Client) GetMovieById(id int) []byte {
	return c.get("/v2/movie/subject/" + strconv.Itoa(id))
}

func (c *Client) GetMoviePhotosById(id int) []byte {
	return c.get("/v2/movie/subject/" + strconv.Itoa(id) + "/photos")
}

//http://developers.douban.com/wiki/?title=movie_v2#reviews
func (c *Client) GetMovieReviewsById(id int) []byte {
	return c.get("/v2/movie/subject/" + strconv.Itoa(id) + "/reviews")
}

//http://developers.douban.com/wiki/?title=movie_v2#comments
func (c *Client) GetMovieCommentsById(id int) []byte {
	return c.get("/v2/movie/subject/" + strconv.Itoa(id) + "/comments")
}

//http://developers.douban.com/wiki/?title=movie_v2#celebrity
func (c *Client) GetCelebrityById(id int) []byte {
	return c.get("/v2/movie/celebrity/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=movie_v2#celebrity-photos
func (c *Client) GetCelebrityPhotosById(id int) []byte {
	return c.get("/v2/movie/celebrity/" + strconv.Itoa(id) + "/photos")
}

//http://developers.douban.com/wiki/?title=movie_v2#works
func (c *Client) GetCelebrityWorksById(id int) []byte {
	return c.get("/v2/movie/celebrity/" + strconv.Itoa(id) + "/works")
}

//http://developers.douban.com/wiki/?title=movie_v2#search
func (c *Client) SearchMovieByKeywords(keywords []string, start, count int) []byte {
	return c.get(Urlencode("/v2/movie/search", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) SearchMovieByTag(tag string, start, count int) []byte {
	return c.get(Urlencode("/v2/movie/search", map[string]string{
		"tag":   tag,
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

//http://developers.douban.com/wiki/?title=movie_v2#nowplaying
func (c *Client) GetNowplayingMovies() []byte {
	return c.get("/v2/movie/nowplaying")
}

//http://developers.douban.com/wiki/?title=movie_v2#coming
func (c *Client) GetComingMovies() []byte {
	return c.get("/v2/movie/coming")
}

//http://developers.douban.com/wiki/?title=movie_v2#top250
func (c *Client) Top250Movies() []byte {
	return c.get("/v2/movie/top250")
}

//http://developers.douban.com/wiki/?title=movie_v2#weekly
func (c *Client) WeeklyMovies() []byte {
	return c.get("/v2/movie/weekly")
}

//http://developers.douban.com/wiki/?title=movie_v2#us-box
func (c *Client) US_Box() []byte {
	return c.get("/v2/movie/us_box")
}

//http://developers.douban.com/wiki/?title=movie_v2#new-movies
func (c *Client) NewMovies() []byte {
	return c.get("/v2/movie/new_movies")
}

//http://developers.douban.com/wiki/?title=music_v2#get_music
func (c *Client) GetMusicById(id int) []byte {
	return c.get("/v2/music/" + strconv.Itoa(id))
}

func (c *Client) SearchMusicByKeywords(keywords []string, start, count int) []byte {
	return c.get(Urlencode("/v2/music/search", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) SearchMusicByTag(tag string, start, count int) []byte {
	return c.get(Urlencode("/v2/music/search", map[string]string{
		"tag":   tag,
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

//http://developers.douban.com/wiki/?title=music_v2#get_music_tags
func (c *Client) GetMusicTagsById(id int) []byte {
	return c.get("/v2/music/" + strconv.Itoa(id) + "/tags")
}

//http://developers.douban.com/wiki/?title=music_v2#post_music_review
func (c *Client) PostMusicReview(id int, title, content string, rating int) []byte {
	return c.post("/v2/music/reviews", map[string]([]string){
		"music":   {strconv.Itoa(id)},
		"title":   {title},
		"content": {content},
		"rating":  {strconv.Itoa(rating)},
	})
}
