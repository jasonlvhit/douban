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

/*
 Basic methods for douban OAuth2 API, get, post, delete and put.
*/

func (c *Client) get(url string) (code int, resp_content []byte) {
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
	code = resp.StatusCode
	resp_content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func (c *Client) post(purl string, data map[string]([]string)) (code int, resp_content []byte) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", API_HOST+purl, strings.NewReader(url.Values(data).Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.access_token["access_token"].(string))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	code = resp.StatusCode
	resp_content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func (c *Client) delete(purl string) (code int, resp_content []byte) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", API_HOST+purl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.access_token["access_token"].(string))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	code = resp.StatusCode
	resp_content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func (c *Client) put(purl string, data map[string]([]string)) (code int, resp_content []byte) {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", API_HOST+purl, strings.NewReader(url.Values(data).Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.access_token["access_token"].(string))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	code = resp.StatusCode
	resp_content, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

/*
 Douban OAuth 2 Authorization implemention.
*/

// Return the authorazation URl.
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

func (c *Client) Me() (code int, resp_content []byte) {
	return c.get("/v2/user/~me")
}

func (c *Client) GetUserbyId(id int) (code int, resp_content []byte) {
	return c.get(fmt.Sprintf("/v2/user/%s", strconv.Itoa(id)))
}

func (c *Client) GetUserbyName(username string) (code int, resp_content []byte) {
	return c.get("/v2/user/" + username)
}

func (c *Client) SearchUserByKeywords(keywords []string, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/user", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) Follow(user_id int) (code int, resp_content []byte) {
	return c.post("/shuo/v2/friendships/create", map[string]([]string){
		"user_id": {strconv.Itoa(user_id)},
	})
}

func (c *Client) Unfollow(user_id int) (code int, resp_content []byte) {
	return c.post("shuo/v2/friendships/destroy", map[string]([]string){
		"user_id": {strconv.Itoa(user_id)},
	})
}

func (c *Client) Following(user_id, start, count int) (code int, resp_content []byte) {
	page := start / count
	return c.get(Urlencode(fmt.Sprintf("shuo/v2/users/%s/following", strconv.Itoa(user_id)), map[string]string{
		"page":  strconv.Itoa(page),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) Followers(user_id, start, count int) (code int, resp_content []byte) {
	page := start / count
	return c.get(Urlencode(fmt.Sprintf("shuo/v2/users/%s/followers", strconv.Itoa(user_id)), map[string]string{
		"page":  strconv.Itoa(page),
		"count": strconv.Itoa(count),
	}))
}

/*
 Book:豆瓣读书
 Url
*/

//http://developers.douban.com/wiki/?title=book_v2#get_book
func (c *Client) GetBookById(id int) (code int, resp_content []byte) {
	return c.get("/v2/book/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=book_v2#get_isbn_book
func (c *Client) GetBookByISBN(isbn string) (code int, resp_content []byte) {
	return c.get("/v2/book/isbn/" + isbn)
}

//http://developers.douban.com/wiki/?title=book_v2#get_book_search
func (c *Client) SearchBookByKeywords(keywords []string, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/book/search", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

//http://developers.douban.com/wiki/?title=book_v2#get_book_search
func (c *Client) SearchBookByTag(tag string, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/book/search", map[string]string{
		"tag":   tag,
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

//http://developers.douban.com/wiki/?title=book_v2#get_book_tags
func (c *Client) GetTagsOfBookById(id int) (code int, resp_content []byte) {
	return c.get("/v2/book/" + strconv.Itoa(id) + "/tags")
}

/*
http://developers.douban.com/wiki/?title=book_v2#get_user_book_tags
func (c *Client) GetTagsOfBooksOfUser() {

}
*/

// http://developers.douban.com/wiki/?title=book_v2#get_user_collections
func (c *Client) GetBookCollectionsOfUser(username string) (code int, resp_content []byte) {
	return c.get("/v2/book/user/" + username + "/collections")
}

// http://developers.douban.com/wiki/?title=book_v2#get_book_collection
func (c *Client) GetBookCollectionOfUser(book_id, username string) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/book/"+book_id+"/collection", map[string]string{
		"user_id": username,
	}))
}

// http://developers.douban.com/wiki/?title=book_v2#get_user_annotations
func (c *Client) GetUserAnnotations(username string) (code int, resp_content []byte) {
	return c.get("/v2/book/user/" + username + "/annotations")
}

// http://developers.douban.com/wiki/?title=book_v2#get_book_annotations
func (c *Client) GetBookAnnotations(book_id int) (code int, resp_content []byte) {
	return c.get("/v2/book/" + strconv.Itoa(book_id) + "/annotations")
}

// http://developers.douban.com/wiki/?title=book_v2#get_annotation
func (c *Client) GetAnnotationById(id int) (code int, resp_content []byte) {
	return c.get("/v2/book/annotation/" + strconv.Itoa(id))
}

// http://developers.douban.com/wiki/?title=book_v2#post_book_collection
// 用户收藏图书
func (c *Client) CollectBookById(id int, status, tags, comment, privacy string, rating int) (code int, resp_content []byte) {
	return c.post(fmt.Sprintf("/v2/book/%s/collection", strconv.Itoa(id)), map[string]([]string){
		"status":  {status},
		"tags":    {tags},
		"comment": {comment},
		"privacy": {privacy},
		"rating":  {strconv.Itoa(rating)},
	})
}

// http://developers.douban.com/wiki/?title=book_v2#put_book_collection
// 用户修改对某本图书的收藏
func (c *Client) UpdateCollectedBookById(id int, status, tags, comment, privacy string, rating int) (code int, resp_content []byte) {
	return c.put(fmt.Sprintf("/v2/book/%s/collection", strconv.Itoa(id)), map[string]([]string){
		"status":  {status},
		"tags":    {tags},
		"comment": {comment},
		"privacy": {privacy},
		"rating":  {strconv.Itoa(rating)},
	})
}

//http://developers.douban.com/wiki/?title=book_v2#delete_book_collection
//用户删除对某本图书的收藏
func (c *Client) DeleteBookCollectionById(id int) (code int, resp_content []byte) {
	return c.delete(fmt.Sprintf("/v2/book/%s/collection", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=book_v2#post_book_annotation
//用户给某本图书写笔记
func (c *Client) MakeAnnotation(id int, content string, page int) (code int, resp_content []byte) {
	return c.post(fmt.Sprintf("/v2/book/%s/annotations", strconv.Itoa(id)), map[string]([]string){
		"content": {content},
		"page":    {strconv.Itoa(page)},
	})
}

//http://developers.douban.com/wiki/?title=book_v2#put_annotation
//用户修改某篇笔记
func (c *Client) UpdateAnnotationById(id int, content string, page int) (code int, resp_content []byte) {
	return c.put(fmt.Sprintf("/v2/book/%s/annotations", strconv.Itoa(id)), map[string]([]string){
		"content": {content},
		"page":    {strconv.Itoa(page)},
	})
}

//http://developers.douban.com/wiki/?title=book_v2#delete_annotation
//用户删除某篇笔记
func (c *Client) DeleteAnnotationById(id int) (code int, resp_content []byte) {
	return c.delete("/v2/book/annotation/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=book_v2#post_book_review
//发表新评论
func (c *Client) NewBookReview(id int, title, content string, rating int) (code int, resp_content []byte) {
	return c.post("/v2/book/reviews", map[string]([]string){
		"book":    {strconv.Itoa(id)},
		"title":   {title},
		"content": {content},
		"rating":  {strconv.Itoa(rating)},
	})
}

//http://developers.douban.com/wiki/?title=book_v2#put_book_review
//修改评论
func (c *Client) UpdateReviewOfBook(id int, title, content string, rating int) (code int, resp_content []byte) {
	return c.put("/v2/book/review/"+strconv.Itoa(id), map[string]([]string){
		"title":   {title},
		"content": {content},
		"rating":  {strconv.Itoa(rating)},
	})
}

//http://developers.douban.com/wiki/?title=book_v2#delete_book_review
//删除评论
func (c *Client) DeleteReviewOfBook(id int) (code int, resp_content []byte) {
	return c.delete("/v2/book/review/" + strconv.Itoa(id))
}

/*
 Movie:豆瓣电影
 Url:
*/

func (c *Client) GetMovieById(id int) (code int, resp_content []byte) {
	return c.get("/v2/movie/subject/" + strconv.Itoa(id))
}

func (c *Client) GetMoviePhotosById(id int) (code int, resp_content []byte) {
	return c.get("/v2/movie/subject/" + strconv.Itoa(id) + "/photos")
}

//http://developers.douban.com/wiki/?title=movie_v2#reviews
func (c *Client) GetMovieReviewsById(id int) (code int, resp_content []byte) {
	return c.get("/v2/movie/subject/" + strconv.Itoa(id) + "/reviews")
}

//http://developers.douban.com/wiki/?title=movie_v2#comments
func (c *Client) GetMovieCommentsById(id int) (code int, resp_content []byte) {
	return c.get("/v2/movie/subject/" + strconv.Itoa(id) + "/comments")
}

//http://developers.douban.com/wiki/?title=movie_v2#celebrity
func (c *Client) GetCelebrityById(id int) (code int, resp_content []byte) {
	return c.get("/v2/movie/celebrity/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=movie_v2#celebrity-photos
func (c *Client) GetCelebrityPhotosById(id int) (code int, resp_content []byte) {
	return c.get("/v2/movie/celebrity/" + strconv.Itoa(id) + "/photos")
}

//http://developers.douban.com/wiki/?title=movie_v2#works
func (c *Client) GetCelebrityWorksById(id int) (code int, resp_content []byte) {
	return c.get("/v2/movie/celebrity/" + strconv.Itoa(id) + "/works")
}

//http://developers.douban.com/wiki/?title=movie_v2#search
func (c *Client) SearchMovieByKeywords(keywords []string, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/movie/search", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) SearchMovieByTag(tag string, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/movie/search", map[string]string{
		"tag":   tag,
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

//http://developers.douban.com/wiki/?title=movie_v2#nowplaying
func (c *Client) GetNowplayingMovies() (code int, resp_content []byte) {
	return c.get("/v2/movie/nowplaying")
}

//http://developers.douban.com/wiki/?title=movie_v2#coming
func (c *Client) GetComingMovies() (code int, resp_content []byte) {
	return c.get("/v2/movie/coming")
}

//http://developers.douban.com/wiki/?title=movie_v2#top250
func (c *Client) Top250Movies() (code int, resp_content []byte) {
	return c.get("/v2/movie/top250")
}

//http://developers.douban.com/wiki/?title=movie_v2#weekly
func (c *Client) WeeklyMovies() (code int, resp_content []byte) {
	return c.get("/v2/movie/weekly")
}

//http://developers.douban.com/wiki/?title=movie_v2#us-box
func (c *Client) US_Box() (code int, resp_content []byte) {
	return c.get("/v2/movie/us_box")
}

//http://developers.douban.com/wiki/?title=movie_v2#new-movies
func (c *Client) NewMovies() (code int, resp_content []byte) {
	return c.get("/v2/movie/new_movies")
}

/*
Music:豆瓣音乐
URL：http://developers.douban.com/wiki/?title=music_v2
*/

//http://developers.douban.com/wiki/?title=music_v2#get_music
func (c *Client) GetMusicById(id int) (code int, resp_content []byte) {
	return c.get("/v2/music/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=music_v2#get_music_search
func (c *Client) SearchMusicByKeywords(keywords []string, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/music/search", map[string]string{
		"q":     strings.Join(keywords, "+"),
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) SearchMusicByTag(tag string, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/music/search", map[string]string{
		"tag":   tag,
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

//http://developers.douban.com/wiki/?title=music_v2#get_music_tags
func (c *Client) GetMusicTagsById(id int) (code int, resp_content []byte) {
	return c.get("/v2/music/" + strconv.Itoa(id) + "/tags")
}

//http://developers.douban.com/wiki/?title=music_v2#post_music_review
//新发音乐评论
func (c *Client) NewMusicReview(id int, title, content string, rating int) (code int, resp_content []byte) {
	return c.post("/v2/music/reviews", map[string]([]string){
		"music":   {strconv.Itoa(id)},
		"title":   {title},
		"content": {content},
		"rating":  {strconv.Itoa(rating)},
	})
}

//http://developers.douban.com/wiki/?title=music_v2#put_music_review
//修改评论
func (c *Client) UpdateMusicReview(id int, title, content string, rating int) (code int, resp_content []byte) {
	return c.put("/v2/music/review/"+strconv.Itoa(id), map[string]([]string){
		"title":   {title},
		"content": {content},
		"rating":  {strconv.Itoa(rating)},
	})
}

//http://developers.douban.com/wiki/?title=music_v2#delete_music_review
//删除评论
func (c *Client) DeleteMusicReview(id int) (code int, resp_content []byte) {
	return c.delete("/v2/music/review/" + strconv.Itoa(id))
}

/*
 photo:豆瓣相册API V2
 URL：http://developers.douban.com/wiki/?title=photo_v2
*/

func (c *Client) GetAlbumById(id int) (code int, resp_content []byte) {
	return c.get("v2/album/" + strconv.Itoa(id))
}

func (c *Client) GetUserCreateAlbums(user_id, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/album/user_created"+strconv.Itoa(user_id), map[string]string{
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) LikeAblumById(id int) (code int, resp_content []byte) {
	return c.post("/v2/album/"+strconv.Itoa(id)+"like", map[string]([]string){})
}

func (c *Client) GetAlbumPhotosById(id, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/album/"+strconv.Itoa(id)+"/photos", map[string]string{
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) GetLikedListByUserId(user_id, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/ablum/user_liked/"+strconv.Itoa(user_id), map[string]string{
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

// http://developers.douban.com/wiki/?title=photo_v2#get_photo
// 获取照片
func (c *Client) GetPhotoById(id int) (code int, resp_content []byte) {
	return c.get("/v2/photo/" + strconv.Itoa(id))
}

/*
 doumail:豆邮API V2
*/

func (c *Client) GetDoumailById(id int) (code int, resp_content []byte) {
	return c.get("/v2/doumail/" + strconv.Itoa(id))
}

func (c *Client) GetInbox(start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/doumail/inbox", map[string]string{
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) GetOutbox(start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/doumail/outbox", map[string]string{
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

func (c *Client) GetUnread(start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode("/v2/doumail/inbox/unread", map[string]string{
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

/*
 online:豆瓣线上活动API
 URL：http://developers.douban.com/wiki/?title=online_v2
*/

//http://developers.douban.com/wiki/?title=online_v2#get
//获取线上活动
func (c *Client) GetOnlineById(id int) (code int, resp_content []byte) {
	return c.get("/v2/online/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=online_v2#users
//获取线上活动成员列表
func (c *Client) GetParticipantsById(id int) (code int, resp_content []byte) {
	return c.get(fmt.Sprintf("/v2/online/%s/participants", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=online_v2#discussions
//获取线上活动论坛列表
func (c *Client) GetOnlineDiscussionById(id int) (code int, resp_content []byte) {
	return c.get(fmt.Sprintf("/v2/online/%s/discussions", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=online_v2#list
//获取线上活动列表
func (c *Client) GetOnlineList() (code int, resp_content []byte) {
	return c.get("/v2/onlines")
}

//http://developers.douban.com/wiki/?title=online_v2#new
//创建活动
func (c *Client) CreateOnline(title, desc, begin_time, end_time,
	related_url, cascade_invite, tags string) (code int, resp_content []byte) {
	return c.post("v2/onlines", map[string]([]string){
		"title":          {title},
		"desc":           {desc},
		"begin_time":     {begin_time},
		"end_time":       {end_time},
		"related_url":    {related_url},
		"cascade_invite": {cascade_invite},
		"tags":           {tags},
	})
}

//http://developers.douban.com/wiki/?title=online_v2#update
//更新线上活动
func (c *Client) UpdateOnline(id int, title, desc, begin_time, end_time,
	related_url, cascade_invite, tags string) (code int, resp_content []byte) {
	return c.put("/v2/online/"+strconv.Itoa(id), map[string]([]string){
		"title":          {title},
		"desc":           {desc},
		"begin_time":     {begin_time},
		"end_time":       {end_time},
		"related_url":    {related_url},
		"cascade_invite": {cascade_invite},
		"tags":           {tags},
	})
}

// http://developers.douban.com/wiki/?title=online_v2#delete
// 删除线上活动
func (c *Client) DeleteOnline(id int) (code int, resp_content []byte) {
	return c.delete("/v2/online/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=online_v2#join
//参加线上活动
func (c *Client) JoinOnline(id int) (code int, resp_content []byte) {
	return c.post(fmt.Sprintf("/v2/online/%s/participants", strconv.Itoa(id)), map[string]([]string){})
}

//http://developers.douban.com/wiki/?title=online_v2#quit
func (c *Client) QuitOnline(id int) (code int, resp_content []byte) {
	return c.delete(fmt.Sprintf("/v2/online/%s/participants", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=online_v2#like
//喜欢线上活动
func (c *Client) LikeOnline(id int) (code int, resp_content []byte) {
	return c.post(fmt.Sprintf("/v2/online/%s/like", strconv.Itoa(id)), map[string]([]string){})
}

//http://developers.douban.com/wiki/?title=online_v2#unlike
//取消喜欢线上活动
func (c *Client) UnlikeOnline(id int) (code int, resp_content []byte) {
	return c.delete(fmt.Sprintf("/v2/online/%s/like", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=online_v2#photos
//线上活动图片列表
func (c *Client) GetOnlinePhotos(id int) (code int, resp_content []byte) {
	return c.get(fmt.Sprintf("/v2/online/%s/photos", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=online_v2#upload
//上传图片
func (c *Client) UploadOnlinePhoto(id int) /*(code int, resp_content []byte)*/ {
	/*
		Album API
		http://developers.douban.com/wiki/?title=photo_v2#new_photo
		Not implemented.
	*/
	//return c.post(purl, data)
}

//http://developers.douban.com/wiki/?title=online_v2#post
//新发讨论
func (c *Client) NewOnlineDiscussion(target_id int, title, content string) (code int, resp_content []byte) {
	return c.NewDiscussion("online", target_id, title, content)
}

//http://developers.douban.com/wiki/?title=online_v2#join_list
//获取用户参加的线上活动列表
func (c *Client) GetUserPaticipateOnline(user_id int) (code int, resp_content []byte) {
	return c.get("/v2/online/user_participated/" + strconv.Itoa(user_id))
}

//http://developers.douban.com/wiki/?title=online_v2#owned
//获取用户创建的线上活动列表
func (c *Client) GetUserCreateOnline(user_id int) (code int, resp_content []byte) {
	return c.get("/v2/online/user_created/" + strconv.Itoa(user_id))
}

/*
 discussion:豆瓣讨论API V2
 URL：http://developers.douban.com/wiki/?title=discussion_v2
*/

// http://developers.douban.com/wiki/?title=discussion_v2#get
// 获取讨论
func (c *Client) GetDiscussionById(id int) (code int, resp_content []byte) {
	return c.get("/v2/discussion/" + strconv.Itoa(id))
}

// http://developers.douban.com/wiki/?title=discussion_v2#update
// 更新讨论
func (c *Client) UpdateDiscussionById(id int, title, content string) (code int, resp_content []byte) {
	return c.put("/v2/discussion/"+strconv.Itoa(id), map[string]([]string){
		"title":   {title},
		"content": {content},
	})
}

// http://developers.douban.com/wiki/?title=discussion_v2#delete
// 删除讨论
func (c *Client) DeleteDiscussionById(id int) (code int, resp_content []byte) {
	return c.delete("/v2/discussion/" + strconv.Itoa(id))
}

// http://developers.douban.com/wiki/?title=discussion_v2#new
// 创建讨论
func (c *Client) NewDiscussion(target string, target_id int, title, content string) (code int, resp_content []byte) {
	return c.post(fmt.Sprintf("/v2/%s/%s/discussions", target, strconv.Itoa(target_id)), map[string]([]string){
		"title":   {title},
		"content": {content},
	})
}

// http://developers.douban.com/wiki/?title=discussion_v2#list
// 获取论坛讨论列表
func (c *Client) GetDiscussionListById(target string, target_id int) (code int, resp_content []byte) {
	return c.get(fmt.Sprintf("/v2/%s/%s/discussions", target, strconv.Itoa(target_id)))

}

/*
 comment:豆瓣回复API V2
 URL：http://developers.douban.com/wiki/?title=comment_v2
*/

// http://developers.douban.com/wiki/?title=comment_v2#list
// 获取回复列表
func (c *Client) GetCommentsById(target string, id, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode(fmt.Sprintf("/v2/%s/%s/comments", target, strconv.Itoa(id)), map[string]string{
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

// http://developers.douban.com/wiki/?title=comment_v2#new
// 新发回复
func (c *Client) NewComment(target string, id int, content string) (code int, resp_content []byte) {
	return c.post(fmt.Sprintf("/v2/%s/%s/comments", target, strconv.Itoa(id)), map[string]([]string){
		"content": {content},
	})
}

// http://developers.douban.com/wiki/?title=comment_v2#get
// 获取单条回复
func (c *Client) GetCommentById(target string, target_id, comment_id int) (code int, resp_content []byte) {
	return c.get(fmt.Sprintf("/v2/%s/%s/comment/%s", target, strconv.Itoa(target_id), strconv.Itoa(comment_id)))
}

// http://developers.douban.com/wiki/?title=comment_v2#delete
// 删除回复
func (c *Client) DeleteCommentById(target string, target_id, comment_id int) (code int, resp_content []byte) {
	return c.delete(fmt.Sprintf("/v2/%s/%s/comment/%s", target, strconv.Itoa(target_id), strconv.Itoa(comment_id)))
}

/*
 travel:豆瓣我去API V2
 URL：http://developers.douban.com/wiki/?title=travel_v2
*/

//http://developers.douban.com/wiki/?title=travel_v2#user-collection
//获取用户收藏地点A
func (c *Client) GetUserCollectLocations(username string, start, count int) (code int, resp_content []byte) {
	return c.get(Urlencode(fmt.Sprintf("/v2/travel/user/%s/collections", username), map[string]string{
		"start": strconv.Itoa(start),
		"count": strconv.Itoa(count),
	}))
}

/*
 event：豆瓣同城
 URL：http://developers.douban.com/wiki/?title=event_v2
*/

//http://developers.douban.com/wiki/?title=event_v2#event_get
//获取活动
func (c *Client) GetEventById(id int) (code int, resp_content []byte) {
	return c.get("/v2/event/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=event_v2#event_participants
//获取参加活动的用户
func (c *Client) GetEventParticipants(id int) (code int, resp_content []byte) {
	return c.get(fmt.Sprintf("/v2/event/%s/participants", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=event_v2#event_wishers
//获取对活动感兴趣的用户
func (c *Client) GetEventWishers(id int) (code int, resp_content []byte) {
	return c.get(fmt.Sprintf("/v2/event/%s/wishers", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=event_v2#event_user_created
//获取用户创建的活动
func (c *Client) GetUserCreateEvents(id int) (code int, resp_content []byte) {
	return c.get("/v2/event/user_created/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=event_v2#event_user_participated
//获取用户参加的活动
func (c *Client) GetUserPartiipatedEvent(id int) (code int, resp_content []byte) {
	return c.get("/v2/event/user_participated/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=event_v2#event_user_wished
//获取用户感兴趣的活动
func (c *Client) GetUserWishedEvents(id int) (code int, resp_content []byte) {
	return c.get("/v2/event/user_wished/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=event_v2#event_list
//获取活动列表
func (c *Client) GetEventList() (code int, resp_content []byte) {
	return c.get("/v2/event/list")
}

//http://developers.douban.com/wiki/?title=event_v2#loc_get
//获取城市
func (c *Client) GetLocationById(id int) (code int, resp_content []byte) {
	return c.get("/v2/loc/" + strconv.Itoa(id))
}

//http://developers.douban.com/wiki/?title=event_v2#loc_list
//获取城市列表
func (c *Client) GetLocationList() (code int, resp_content []byte) {
	return c.get("/v2/loc/list")
}

//http://developers.douban.com/wiki/?title=event_v2#event_participate
// 参加活动
func (c *Client) ParticipateEventById(id int) (code int, resp_content []byte) {
	return c.post(fmt.Sprintf("/v2/event/%s/participants", strconv.Itoa(id)), map[string]([]string){})
}

//http://developers.douban.com/wiki/?title=event_v2#event_quit
//不参加活动
func (c *Client) QuitEventById(id int) (code int, resp_content []byte) {
	return c.delete(fmt.Sprintf("/v2/event/%s/participants", strconv.Itoa(id)))
}

//http://developers.douban.com/wiki/?title=event_v2#event_wish
//对活动感兴趣
func (c *Client) WishEvent(id int) (code int, resp_content []byte) {
	return c.post(fmt.Sprintf("/v2/event/%s/wishers", strconv.Itoa(id)), map[string]([]string){})
}

//http://developers.douban.com/wiki/?title=event_v2#event_unwish
//不感兴趣
func (c *Client) UnwishEvent(id int) (code int, resp_content []byte) {
	return c.delete(fmt.Sprintf("/v2/event/%s/wishers", strconv.Itoa(id)))
}
