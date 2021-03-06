package myanimelist

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Forum struct {
	mal *MAL
}

// ForumBoards return list of all forum's categories.
func (f *Forum) Boards() (*ForumCategories, error) {
	method := http.MethodGet
	path := "./forum/boards"

	data := url.Values{}

	var categories = &ForumCategories{}
	if err := f.mal.request(categories, method, path, data); err != nil {
		return nil, err
	}

	return categories, nil
}

// ForumCategories stores data received from executing ForumBoards()
type ForumCategories struct {
	Categories []ForumCategory `json:"categories"`
}
type ForumCategory struct {
	Title  string       `json:"title"`
	Boards []ForumBoard `json:"boards"`
}
type ForumBoard struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	SubBoards   []ForumSubBoard `json:"subboards"`
}
type ForumSubBoard struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// ForumTopic retrieves info about topic with provided topicID.
func (f *Forum) Topic(topicID int, settings PagingSettings) (*ForumTopic, error) {
	method := http.MethodGet
	path := fmt.Sprintf("./forum/topic/%d", topicID)

	data := url.Values{}

	settings.set(&data)

	topicInfo := &ForumTopic{parent: f}

	if err := f.mal.request(topicInfo, method, path, data); err != nil {
		return nil, err
	}
	return topicInfo, nil
}

// ForumTopic stores topic title, poll, array of posts.
// Use Prev() and Next() methods to retrieve corresponding result pages.
type ForumTopic struct {
	parent *Forum
	Data   struct {
		Title string      `json:"title"`
		Posts []ForumPost `json:"posts"`
		Poll  Poll        `json:"poll"`
	} `json:"data"`
	Paging Paging `json:"paging"`
}
type ForumPost struct {
	ID        int                 `json:"id"`
	Number    int                 `json:"number"`
	CreatedAt time.Time           `json:"created_at"`
	CreatedBy ForumUserWithAvatar `json:"created_by"`
	Body      string              `json:"body"`
	Signature string              `json:"signature"`
}
type ForumUserWithAvatar struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"forum_avator"`
}
type Poll struct {
	ID       int          `json:"id"`
	Question string       `json:"question"`
	Closed   bool         `json:"closed"`
	Options  []PollOption `json:"options"`
}
type PollOption struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *ForumTopic) Prev(limit ...int) (result *ForumTopic, err error) {
	result = &ForumTopic{parent: obj.parent}
	err = obj.parent.mal.getPage(result, obj.Paging, -1, limit)
	return
}

// Next return next result page.
// If its last page - returns error.
func (obj *ForumTopic) Next(limit ...int) (result *ForumTopic, err error) {
	result = &ForumTopic{parent: obj.parent}
	err = obj.parent.mal.getPage(result, obj.Paging, 1, limit)
	return
}

// ForumSearchSetting represent advanced search on MyAnimeList forum.
// All fields are optional.
type ForumSearchSettings struct {
	Keyword      string
	BoardID      int
	SubboardID   int
	TopicStarter string
	PostAuthor   string
}

// ForumSearchTopics implements advanced search from website.
// Use ForumSearchSettings struct to set search options.
func (f *Forum) Search(searchOpts ForumSearchSettings, settings PagingSettings) (*ForumSearchResult, error) {
	method := http.MethodGet
	path := "./forum/topics"

	data := url.Values{}

	// only this sort method working yet
	sort := "recent"

	if searchOpts.Keyword != "" {
		data.Set("q", searchOpts.Keyword)
	}
	if searchOpts.BoardID != 0 {
		data.Set("board_id", strconv.Itoa(searchOpts.BoardID))
	}
	if searchOpts.SubboardID != 0 {
		data.Set("subboard_id", strconv.Itoa(searchOpts.SubboardID))
	}
	if searchOpts.TopicStarter != "" {
		data.Set("topic_user_name", searchOpts.TopicStarter)
	}
	if searchOpts.PostAuthor != "" {
		data.Set("user_name", searchOpts.PostAuthor)
	}

	data.Set("sort", sort)

	settings.set(&data)

	result := &ForumSearchResult{parent: f}
	if err := f.mal.request(result, method, path, data); err != nil {
		return nil, err
	}

	return result, nil
}

// ForumSearchResult stores array with search result entries.
// Use Prev() and Next() methods to retrieve corresponding result pages.
type ForumSearchResult struct {
	parent *Forum
	Data   []ForumSearchEntry `json:"data"`
	Paging Paging             `json:"paging"`
}
type ForumSearchEntry struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	CreatedAt         time.Time `json:"created_at"`
	CreatedBy         ForumUser `json:"created_by"`
	NumberOfPosts     int       `json:"number_of_posts"`
	LastPostCreatedAt time.Time `json:"last_post_created_at"`
	LastPostCreatedBy ForumUser `json:"last_post_created_by"`
	IsLocked          bool      `json:"is_locked"`
}
type ForumUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *ForumSearchResult) Prev(limit ...int) (result *ForumSearchResult, err error) {
	result = &ForumSearchResult{parent: obj.parent}
	err = obj.parent.mal.getPage(result, obj.Paging, -1, limit)
	return
}

// Next return next result page.
// If its last page - returns error.
func (obj *ForumSearchResult) Next(limit ...int) (result *ForumSearchResult, err error) {
	result = &ForumSearchResult{parent: obj.parent}
	err = obj.parent.mal.getPage(result, obj.Paging, 1, limit)
	return
}
