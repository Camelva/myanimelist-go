package myanimelist

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DeleteAnimeFromList remove entry with certain ID from current user's list.
// If the specified entry does not exist in user's list,
// function acts like call was successful and returns nil
func (mal *MAL) DeleteAnimeFromList(animeID int) error {
	method := http.MethodDelete
	path := fmt.Sprintf("./anime/%d/my_list_status", animeID)
	query := url.Values{}

	resp, err := mal.requestRaw(method, path, query)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

// UpdateAnimeStatus changes specified anime' properties according to provided AnimeConfig.
// Returns updated AnimeStatus or error, if any.
func (mal *MAL) UpdateAnimeStatus(config AnimeConfig) (*AnimeStatus, error) {
	method := http.MethodPatch

	animeID, ok := config["id"]
	if !ok {
		return nil, errors.New("anime ID is required")
	}
	delete(config, "id")

	path := fmt.Sprintf("./anime/%v/my_list_status", animeID)
	data := url.Values{}
	for k, v := range config {
		data.Set(k, v)
	}

	animeS := new(AnimeStatus)
	if err := mal.request(animeS, method, path, data); err != nil {
		return nil, err
	}

	return animeS, nil
}

// NewAnimeConfig generates new config and sets anime id
// You can also manually create object and set id with SetID(id int) method
func NewAnimeConfig(id int) AnimeConfig {
	return AnimeConfig{"id": strconv.Itoa(id)}
}

// AnimeConfig contains all the changes you want to apply
type AnimeConfig map[string]string

func (c AnimeConfig) SetID(id int) AnimeConfig {
	c["id"] = strconv.Itoa(id)
	return c
}

// SetStatus accept only StatusWatching, StatusCompleted, StatusOnHold,
// StatusDropped or StatusPlanToWatch constants
func (c AnimeConfig) SetStatus(status string) AnimeConfig {
	acceptableStatuses := makeList(append(GeneralStatuses, AnimeStatuses...))
	if _, ok := acceptableStatuses[status]; !ok {
		// non-acceptable status, do nothing
		return c
	}
	c["status"] = status
	return c
}

func (c AnimeConfig) SetIsRewatching(b bool) AnimeConfig {
	c["is_rewatching"] = fmt.Sprintf("%t", b)
	return c
}

func (c AnimeConfig) SetScore(score int) AnimeConfig {
	if score < 0 || score > 10 {
		return c
	}
	c["score"] = strconv.Itoa(score)
	return c
}

func (c AnimeConfig) SetWatchedEpisodes(count int) AnimeConfig {
	c["num_watched_episodes"] = strconv.Itoa(count)
	return c
}

// SetPriority accept only PriorityLow, PriorityMedium or PriorityHigh constants
func (c AnimeConfig) SetPriority(priority int) AnimeConfig {
	acceptable := makeListInt(Priorities)
	if _, ok := acceptable[priority]; !ok {
		// non-acceptable, do nothing
		return c
	}
	c["priority"] = strconv.Itoa(int(priority))
	return c
}

func (c AnimeConfig) SetRewatchedTimes(count int) AnimeConfig {
	c["num_times_rewatched"] = strconv.Itoa(count)
	return c
}

func (c AnimeConfig) SetRewatchValue(value int) AnimeConfig {
	if value < 0 || value > 5 {
		return c
	}
	c["rewatch_value"] = strconv.Itoa(value)
	return c
}

func (c AnimeConfig) SetTags(tags ...string) AnimeConfig {
	text := strings.Join(tags, ", ")
	c["tags"] = text
	return c
}

func (c AnimeConfig) SetComment(text string) AnimeConfig {
	c["comments"] = text
	return c
}

// AnimeStatus contains server response about certain anime
type AnimeStatus struct {
	Status             string    `json:"status"`
	Score              int       `json:"score"`
	NumWatchedEpisodes int       `json:"num_episodes_watched"`
	IsRewatching       bool      `json:"is_rewatching"`
	UpdatedAt          time.Time `json:"updated_at"`
	Priority           int       `json:"priority"`
	NumTimesRewatched  int       `json:"num_times_rewatched"`
	RewatchValue       int       `json:"rewatch_value"`
	Tags               []string  `json:"tags"`
	Comments           string    `json:"comments"`
}

// UserAnimeList returns anime list of certain user with provided username (for current user use empty string).
// You can set status to retrieve only anime's with same status or use empty object.
// You can sort list by using on of these constants: SortListByScore, SortListByUpdateDate,
// SortListByTitle, SortListByStartDate, SortListByID or provide empty object to disable sorting
func (mal *MAL) UserAnimeList(username string, status string, sort string, settings PagingSettings) (*UserAnimeList, error) {
	if username == "" {
		username = "@me"
	}

	method := http.MethodGet
	path := fmt.Sprintf("./users/%s/animelist", username)
	data := url.Values{}
	if status != "" {
		acceptable := makeList(append(GeneralStatuses, AnimeStatuses...))
		if _, ok := acceptable[sort]; ok {
			data.Add("status", string(status))
		}
	}
	if sort != "" {
		acceptable := makeList(ListSortings)
		if _, ok := acceptable[sort]; ok {
			data.Add("sort", fixSorting(sort, "anime"))
		}
	}

	settings.Set(&data)

	var userList = new(UserAnimeList)
	if err := mal.request(userList, method, path, data); err != nil {
		return nil, err
	}

	return userList, nil
}

type UserAnimeList struct {
	Data []struct {
		Node       `json:"node"`
		ListStatus AnimeListStatus `json:"list_status"`
	} `json:"data"`
	Paging Paging `json:"paging"`
}

// Next return next result page.
// If its last page - returns error.
func (obj *UserAnimeList) Next(mal *MAL) (*UserAnimeList, error) {
	nextPageURL := obj.Paging.Next
	if nextPageURL == "" {
		return nil, errors.New("latest page")
	}
	nextPageData := new(UserAnimeList)
	if err := mal.request(nextPageData, http.MethodGet, nextPageURL, url.Values{}); err != nil {
		return nil, err
	}
	return nextPageData, nil
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *UserAnimeList) Prev(mal *MAL) (*UserAnimeList, error) {
	prevPageURL := obj.Paging.Previous
	if prevPageURL == "" {
		return nil, errors.New("first page")
	}
	prevPageData := new(UserAnimeList)
	if err := mal.request(prevPageData, http.MethodGet, prevPageURL, url.Values{}); err != nil {
		return nil, err
	}
	return prevPageData, nil
}

type AnimeListStatus struct {
	Status             string    `json:"status"`
	Score              int       `json:"score"`
	NumWatchedEpisodes int       `json:"num_watched_episodes"`
	IsRewatching       bool      `json:"is_rewatching"`
	UpdatedAt          time.Time `json:"updated_at"`
}
