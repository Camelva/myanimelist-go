package myanimelist

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// MangaSearch return list of manga, performing search for similar as provided search string.
func (mal *MAL) MangaSearch(search string, settings PagingSettings) (*MangaSearchResult, error) {
	method := http.MethodGet
	path := "./manga"
	data := url.Values{
		"q": {search},
	}
	settings.Set(&data)

	searchResult := new(MangaSearchResult)
	if err := mal.request(searchResult, method, path, data); err != nil {
		return nil, err
	}
	return searchResult, nil
}

// MangaSearchResult stores array with search entries.
// Use Prev() and Next() methods to retrieve corresponding result pages.
type MangaSearchResult struct {
	Data []struct {
		Node `json:"node"`
	} `json:"data"`
	Paging Paging `json:"paging"`
}

// Next return next result page.
// If its last page - returns error.
func (obj *MangaSearchResult) Next(mal *MAL) (*MangaSearchResult, error) {
	nextPageURL := obj.Paging.Next
	if nextPageURL == "" {
		return nil, errors.New("latest page")
	}
	nextPageData := new(MangaSearchResult)
	if err := mal.request(nextPageData, http.MethodGet, nextPageURL, url.Values{}); err != nil {
		return nil, err
	}
	return nextPageData, nil
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *MangaSearchResult) Prev(mal *MAL) (*MangaSearchResult, error) {
	prevPageURL := obj.Paging.Previous
	if prevPageURL == "" {
		return nil, errors.New("first page")
	}
	prevPageData := new(MangaSearchResult)
	if err := mal.request(prevPageData, http.MethodGet, prevPageURL, url.Values{}); err != nil {
		return nil, err
	}
	return prevPageData, nil
}

// MangaDetails returns details about manga with provided ID.
// You can control which fields to retrieve. For all fields use FieldAllAvailable.
// With no fields provided api still returns ID, Title and MainPicture fields
func (mal *MAL) MangaDetails(mangaID int, fields ...string) (*MangaDetails, error) {
	acceptableArr := append(GeneralFields, MangaFields...)
	method := http.MethodGet
	path := fmt.Sprintf("./manga/%d", mangaID)
	acceptable := makeList(acceptableArr)

	if fields[0] == FieldAllAvailable {
		fields = acceptableArr
	}

	data := url.Values{}
	fieldsString := ""
	for k, f := range fields {
		if _, ok := acceptable[f]; !ok {
			continue
		}
		if k != 0 {
			fieldsString += ", "
		}
		fieldsString += string(f)
	}
	data.Set("fields", fieldsString)

	manga := new(MangaDetails)
	if err := mal.request(manga, method, path, data); err != nil {
		return nil, err
	}

	return manga, nil
}

// MangaDetails contain info about certain manga.
type MangaDetails struct {
	ID                int     `json:"id"`
	Title             string  `json:"title"`
	MainPicture       Picture `json:"main_picture"`
	AlternativeTitles struct {
		Synonyms []string `json:"synonyms"`
		En       string   `json:"en"`
		Ja       string   `json:"ja"`
	} `json:"alternative_titles"`
	StartDate       string    `json:"start_date"`
	EndDate         string    `json:"end_date"`
	Synopsis        string    `json:"synopsis"`
	Mean            float64   `json:"mean"`
	Rank            int       `json:"rank"`
	Popularity      int       `json:"popularity"`
	NumListUsers    int       `json:"num_list_users"`
	NumScoringUsers int       `json:"num_scoring_users"`
	Nsfw            string    `json:"nsfw"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	MediaType       string    `json:"media_type"`
	Status          string    `json:"status"`
	Genres          []Genre   `json:"genres"`
	MyListStatus    struct {
		Status          string    `json:"status"`
		IsRereading     bool      `json:"is_rereading"`
		NumVolumesRead  int       `json:"num_volumes_read"`
		NumChaptersRead int       `json:"num_chapters_read"`
		Score           int       `json:"score"`
		UpdatedAt       time.Time `json:"updated_at"`
	} `json:"my_list_status"`
	NumVolumes  int `json:"num_volumes"`
	NumChapters int `json:"num_chapters"`
	Authors     []struct {
		Node struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		} `json:"node"`
		Role string `json:"role"`
	} `json:"authors"`
	Pictures     []Picture `json:"pictures"`
	Background   string    `json:"background"`
	RelatedAnime []struct {
		Node                  `json:"node"`
		RelationType          string `json:"relation_type"`
		RelationTypeFormatted string `json:"relation_type_formatted"`
	} `json:"related_anime"`
	RelatedManga []struct {
		Node                  `json:"node"`
		RelationType          string `json:"relation_type"`
		RelationTypeFormatted string `json:"relation_type_formatted"`
	} `json:"related_manga"`
	Recommendations []struct {
		Node               `json:"node"`
		NumRecommendations int `json:"num_recommendations"`
	} `json:"recommendations"`
	Serialization []struct {
		Node struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"node"`
	} `json:"serialization"`
}

// MangaRanking returns list of top manga, for each measurement.
// For additional info, see: https://myanimelist.net/apiconfig/references/api/v2#operation/manga_ranking_get
// Currently available ranks:
// - RankAll, - RankManga, - RankNovels, - RankOneShots, - RankDoujinshi,
// - RankManhwa, - RankManhua, - RankByPopularity, - RankFavorite.
func (mal *MAL) MangaRanking(rankingType string, settings PagingSettings) (*MangaRanking, error) {
	// Current working rankings
	acceptable := makeList(append(GeneralRankings, MangaRankings...))
	if _, ok := acceptable[rankingType]; !ok {
		return nil, errors.New("undefined ranking type")
	}

	method := http.MethodGet
	path := "./manga/ranking"
	data := url.Values{
		"ranking_type": {string(rankingType)},
	}
	settings.Set(&data)

	mangaRank := new(MangaRanking)
	if err := mal.request(mangaRank, method, path, data); err != nil {
		return nil, err
	}

	return mangaRank, nil
}

// MangaRanking contain arrays of Nodes (ID, Title, MainPicture) with their rank position.
// Use Prev() and Next() methods to retrieve corresponding result pages.
type MangaRanking struct {
	Data []struct {
		Node    `json:"node"`
		Ranking struct {
			Rank int `json:"rank"`
		} `json:"ranking"`
	} `json:"data"`
	Paging Paging `json:"paging"`
}

// Next return next result page.
// If its last page - returns error.
func (obj *MangaRanking) Next(mal *MAL) (*MangaRanking, error) {
	nextPageURL := obj.Paging.Next
	if nextPageURL == "" {
		return nil, errors.New("latest page")
	}
	nextPageData := new(MangaRanking)
	if err := mal.request(nextPageData, http.MethodGet, nextPageURL, url.Values{}); err != nil {
		return nil, err
	}
	return nextPageData, nil
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *MangaRanking) Prev(mal *MAL) (*MangaRanking, error) {
	prevPageURL := obj.Paging.Previous
	if prevPageURL == "" {
		return nil, errors.New("first page")
	}
	prevPageData := new(MangaRanking)
	if err := mal.request(prevPageData, http.MethodGet, prevPageURL, url.Values{}); err != nil {
		return nil, err
	}
	return prevPageData, nil
}
