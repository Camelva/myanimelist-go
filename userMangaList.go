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

type MangaList struct {
	manga *Manga
}

// DeleteMangaFromList remove entry with certain ID from current user's list.
// If the specified entry does not exist in user's list,
// function acts like call was successful and returns nil
func (ml *MangaList) Remove(animeID int) error {
	method := http.MethodDelete
	path := fmt.Sprintf("./manga/%d/my_list_status", animeID)

	query := url.Values{}

	resp, err := ml.manga.mal.requestRaw(method, path, query)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

// NewMangaConfig generates new config and sets manga id
// You can also manually create object and set id with SetID(id int) method
func NewMangaConfig(id int) MangaConfig {
	return MangaConfig{"id": strconv.Itoa(id)}
}

// MangaConfig contains all the changes you want to apply
type MangaConfig map[string]string

func (c MangaConfig) SetID(id int) MangaConfig {
	c["id"] = strconv.Itoa(id)
	return c
}

// SetStatus accept only StatusReading, StatusCompleted, StatusOnHold,
// StatusDropped or StatusPlanToRead constants
func (c MangaConfig) SetStatus(status string) MangaConfig {
	acceptable := makeList(append(generalStatuses, mangaStatuses...))
	if _, ok := acceptable[status]; !ok {
		// non-acceptable status, do nothing
		return c
	}
	c["status"] = status
	return c
}
func (c MangaConfig) SetIsRereading(b bool) MangaConfig {
	c["is_rereading"] = fmt.Sprintf("%t", b)
	return c
}
func (c MangaConfig) SetScore(score int) MangaConfig {
	if score < 0 || score > 10 {
		return c
	}
	c["score"] = strconv.Itoa(score)
	return c
}
func (c MangaConfig) SetVolumesRead(num int) MangaConfig {
	c["num_volumes_read"] = strconv.Itoa(num)
	return c
}
func (c MangaConfig) SetChaptersRead(num int) MangaConfig {
	c["num_chapters_read"] = strconv.Itoa(num)
	return c
}

// SetPriority accept only PriorityLow, PriorityMedium or PriorityHigh constants
func (c MangaConfig) SetPriority(priority int) MangaConfig {
	acceptable := makeListInt(priorities)
	if _, ok := acceptable[priority]; !ok {
		// non-acceptable, do nothing
		return c
	}
	c["priority"] = strconv.Itoa(priority)
	return c
}
func (c MangaConfig) SetRereadTimes(count int) MangaConfig {
	c["num_times_reread"] = strconv.Itoa(count)
	return c
}
func (c MangaConfig) SetRereadValue(value int) MangaConfig {
	if value < 0 || value > 5 {
		return c
	}
	c["reread_value"] = strconv.Itoa(value)
	return c
}
func (c MangaConfig) SetTags(tags ...string) MangaConfig {
	text := strings.Join(tags, ", ")
	c["tags"] = text
	return c
}
func (c MangaConfig) SetComment(text string) MangaConfig {
	c["comments"] = text
	return c
}

// UpdateMangaStatus changes specified manga' properties according to provided MangaConfig.
// Returns updated MangaStatus or error, if any.
func (ml *MangaList) Update(config MangaConfig) (*MangaStatus, error) {
	method := http.MethodPatch

	mangaID, ok := config["id"]
	if !ok {
		return nil, errors.New("manga ID is required")
	}
	delete(config, "id")

	path := fmt.Sprintf("./manga/%v/my_list_status", mangaID)
	data := url.Values{}
	for k, v := range config {
		data.Set(k, v)
	}

	mangaS := new(MangaStatus)
	if err := ml.manga.mal.request(mangaS, method, path, data); err != nil {
		return nil, err
	}

	return mangaS, nil
}

// MangaStatus contains server response about certain manga
type MangaStatus struct {
	Status          string    `json:"status"`
	IsRereading     bool      `json:"is_rereading"`
	NumVolumesRead  int       `json:"num_volumes_read"`
	NumChaptersRead int       `json:"num_chapters_read"`
	Score           int       `json:"score"`
	UpdatedAt       time.Time `json:"updated_at"`
	Priority        int       `json:"priority"`
	NumTimesReread  int       `json:"num_times_reread"`
	RereadValue     int       `json:"reread_value"`
	Tags            []string  `json:"tags"`
	Comments        string    `json:"comments"`
}

// UserMangaList returns manga list of certain user with provided username
// You can set status to retrieve only manga's with same status or use empty object
// You can sort list by using on of these constants: SortListByScore, SortListByUpdateDate,
// SortListByTitle, SortListByStartDate, SortListByID or provide empty object to disable sorting
func (ml *MangaList) User(username string, status string, sort string, settings PagingSettings) (*UserMangaList, error) {
	if username == "" {
		username = "@me"
	}

	method := http.MethodGet
	path := fmt.Sprintf("./users/%s/mangalist", username)
	data := url.Values{}
	if status != "" {
		acceptable := makeList(append(generalStatuses, mangaStatuses...))
		if _, ok := acceptable[status]; ok {
			data.Add("status", status)
		}
	}
	if sort != "" {
		acceptable := makeList(listSortings)
		if _, ok := acceptable[sort]; ok {
			data.Add("sort", fixSorting(sort, "manga"))
		}
	}
	settings.set(&data)

	var userList = &UserMangaList{parent: ml}
	if err := ml.manga.mal.request(userList, method, path, data); err != nil {
		return nil, err
	}

	return userList, nil
}

type UserMangaList struct {
	parent *MangaList
	Data   []struct {
		Node       `json:"node"`
		ListStatus MangaListStatus `json:"list_status"`
	} `json:"data"`
	Paging Paging `json:"paging"`
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *UserMangaList) Prev(limit ...int) (result *UserMangaList, err error) {
	result = &UserMangaList{parent: obj.parent}
	err = obj.parent.manga.mal.getPage(result, obj.Paging, -1, limit)
	return
}

// Next return next result page.
// If its last page - returns error.
func (obj *UserMangaList) Next(limit ...int) (result *UserMangaList, err error) {
	result = &UserMangaList{parent: obj.parent}
	err = obj.parent.manga.mal.getPage(result, obj.Paging, 1, limit)
	return
}

type MangaListStatus struct {
	Status          string    `json:"status"`
	Score           int       `json:"score"`
	NumVolumesRead  int       `json:"num_volumes_read"`
	NumChaptersRead int       `json:"num_chapters_read"`
	IsRereading     bool      `json:"is_rereading"`
	UpdatedAt       time.Time `json:"updated_at"`
}
