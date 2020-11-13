package myanimelist

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// AnimeSearch return list of anime, performing search for similar as provided search string.
func (mal *MAL) AnimeSearch(search string, settings PagingSettings) (*AnimeSearchResult, error) {
	method := http.MethodGet
	path := "./anime"
	data := url.Values{
		"q": {search},
	}
	settings.set(&data)

	searchResult := new(AnimeSearchResult)
	if err := mal.request(searchResult, method, path, data); err != nil {
		return nil, err
	}
	return searchResult, nil
}

// AnimeSearchResult stores array with search entries.
// Use Prev() and Next() methods to retrieve corresponding result pages.
type AnimeSearchResult struct {
	Data []struct {
		Node `json:"node"`
	} `json:"data"`
	Paging Paging `json:"paging"`
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *AnimeSearchResult) Prev(client *MAL, limit ...int) (result *AnimeSearchResult, err error) {
	result = new(AnimeSearchResult)
	err = client.getPage(result, obj.Paging, -1, limit)
	return
}

// Next return next result page.
// If its last page - returns error.
func (obj *AnimeSearchResult) Next(client *MAL, limit ...int) (result *AnimeSearchResult, err error) {
	result = new(AnimeSearchResult)
	err = client.getPage(result, obj.Paging, 1, limit)
	return
}

// AnimeDetails returns details about anime with provided ID.
// You can control which fields to retrieve. For all fields use FieldAllAvailable.
// With no fields provided api still returns ID, Title and MainPicture fields
func (mal *MAL) AnimeDetails(animeID int, fields ...string) (*AnimeDetails, error) {
	var acceptableArr = append(generalFields, animeFields...)
	method := http.MethodGet
	path := fmt.Sprintf("./anime/%d", animeID)
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
		fieldsString += f
	}
	data.Set("fields", fieldsString)

	anime := new(AnimeDetails)
	if err := mal.request(anime, method, path, data); err != nil {
		return nil, err
	}

	return anime, nil
}

// AnimeDetails contain info about certain anime.
type AnimeDetails struct {
	ID                int     `json:"id"`
	Title             string  `json:"title"`
	MainPicture       Picture `json:"main_picture"`
	AlternativeTitles struct {
		Synonyms []string `json:"synonyms"`
		En       string   `json:"en"`
		Ja       string   `json:"ja"`
	} `json:"alternative_titles"`
	StartDate       string          `json:"start_date"`
	EndDate         string          `json:"end_date"`
	Synopsis        string          `json:"synopsis"`
	Mean            float64         `json:"mean"`
	Rank            int             `json:"rank"`
	Popularity      int             `json:"popularity"`
	NumListUsers    int             `json:"num_list_users"`
	NumScoringUsers int             `json:"num_scoring_users"`
	Nsfw            string          `json:"nsfw"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	MediaType       string          `json:"media_type"`
	Status          string          `json:"status"`
	Genres          []Genre         `json:"genres"`
	MyListStatus    AnimeListStatus `json:"my_list_status"`
	NumEpisodes     int             `json:"num_episodes"`
	StartSeason     struct {
		Year   int    `json:"year"`
		Season string `json:"season"`
	} `json:"start_season"`
	Broadcast struct {
		DayOfTheWeek string `json:"day_of_the_week"`
		StartTime    string `json:"start_time"`
	} `json:"broadcast"`
	Source                 string    `json:"source"`
	AverageEpisodeDuration int       `json:"average_episode_duration"`
	Rating                 string    `json:"rating"`
	Pictures               []Picture `json:"pictures"`
	Background             string    `json:"background"`
	RelatedAnime           []struct {
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
	Studios []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"studios"`
	Statistics struct {
		Status struct {
			Watching    string `json:"watching"`
			Completed   string `json:"completed"`
			OnHold      string `json:"on_hold"`
			Dropped     string `json:"dropped"`
			PlanToWatch string `json:"plan_to_watch"`
		} `json:"status"`
		NumListUsers int `json:"num_list_users"`
	} `json:"statistics"`
}

// AnimeRanking returns list of top anime, for each measurement.
// For additional info, see: https://myanimelist.net/apiconfig/references/api/v2#operation/anime_ranking_get
// Currently available ranks:
// - RankAll, - RankAiring, - RankUpcoming, - RankTV, - RankOVA,
// - RankMovie, - RankSpecial, - RankByPopularity, - RankFavorite.
func (mal *MAL) AnimeRanking(rankingType string, settings PagingSettings) (*AnimeRanking, error) {
	// Currently working rankings
	acceptable := makeList(append(generalRankings, animeRankings...))
	if _, ok := acceptable[rankingType]; !ok {
		return nil, errors.New("undefined ranking type")
	}

	method := http.MethodGet
	path := "./anime/ranking"
	data := url.Values{
		"ranking_type": {rankingType},
	}
	settings.set(&data)

	animeRank := new(AnimeRanking)
	if err := mal.request(animeRank, method, path, data); err != nil {
		return nil, err
	}

	return animeRank, nil
}

// AnimeRanking contain arrays of Nodes (ID, Title, MainPicture) with their rank position.
// Use Prev() and Next() methods to retrieve corresponding result pages.
type AnimeRanking struct {
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
func (obj *AnimeRanking) Next(client *MAL, limit ...int) (result *AnimeRanking, err error) {
	result = new(AnimeRanking)
	err = client.getPage(result, obj.Paging, 1, limit)
	return
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *AnimeRanking) Prev(client *MAL, limit ...int) (result *AnimeRanking, err error) {
	result = new(AnimeRanking)
	err = client.getPage(result, obj.Paging, -1, limit)
	return
}

// SeasonalAnime returns list of anime from certain year's season.
// Year and Season fields are required. Rest are optional.
// For additional info see https://myanimelist.net/apiconfig/references/api/v2#operation/anime_ranking_get
func (mal *MAL) SeasonalAnime(year int, season string, sort string, settings PagingSettings) (*SeasonalAnime, error) {
	// Available season values
	acceptable := makeList(seasons)
	if _, ok := acceptable[season]; !ok {
		return nil, errors.New("undefined season")
	}
	// Available sort values
	availableSorting := map[string]struct{}{SortByScore: {}, SortByUsersLists: {}}

	if year <= 0 {
		return nil, errors.New("invalid year")
	}

	method := http.MethodGet
	path := fmt.Sprintf("./anime/season/%d/%s", year, season)
	data := url.Values{}
	if sort != "" {
		if _, ok := availableSorting[sort]; ok {
			data.Set("sort", sort)
		}
	}
	settings.set(&data)

	seasonal := new(SeasonalAnime)
	if err := mal.request(seasonal, method, path, data); err != nil {
		return nil, err
	}

	return seasonal, nil
}

// SeasonalAnime contain array with basic anime nodes (ID, Title, MainPicture).
// Use Prev() and Next() methods to retrieve corresponding result pages.
type SeasonalAnime struct {
	Data []struct {
		Node `json:"node"`
	} `json:"data"`
	Paging Paging `json:"paging"`
	Season struct {
		Year   int    `json:"year"`
		Season string `json:"season"`
	} `json:"season"`
}

// Next return next result page.
// If its last page - returns error.
func (obj *SeasonalAnime) Next(client *MAL, limit ...int) (result *SeasonalAnime, err error) {
	result = new(SeasonalAnime)
	err = client.getPage(result, obj.Paging, 1, limit)
	return
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *SeasonalAnime) Prev(client *MAL, limit ...int) (result *SeasonalAnime, err error) {
	result = new(SeasonalAnime)
	err = client.getPage(result, obj.Paging, -1, limit)
	return
}

// SuggestedAnime returns suggested anime for the authorized user.
// If the user is new comer, expect to receive empty result.
func (mal *MAL) SuggestedAnime(settings PagingSettings) (*SuggestedAnime, error) {
	method := http.MethodGet
	path := "./anime/suggestions"
	data := url.Values{}
	settings.set(&data)

	suggestions := new(SuggestedAnime)
	if err := mal.request(suggestions, method, path, data); err != nil {
		return nil, err
	}

	return suggestions, nil
}

// SuggestedAnime contain arrays of anime Nodes (ID, Title, MainPicture), suggested for current user.
// Use Prev() and Next() methods to retrieve corresponding result pages.
type SuggestedAnime struct {
	Data []struct {
		Node `json:"node"`
	} `json:"data"`
	Paging Paging `json:"paging"`
}

// Prev return previous result page.
// If its first page - returns error.
func (obj *SuggestedAnime) Prev(client *MAL, limit ...int) (result *SuggestedAnime, err error) {
	result = new(SuggestedAnime)
	err = client.getPage(result, obj.Paging, -1, limit)
	return
}

// Next return next result page.
// If its last page - returns error.
func (obj *SuggestedAnime) Next(client *MAL, limit ...int) (result *SuggestedAnime, err error) {
	result = new(SuggestedAnime)
	err = client.getPage(result, obj.Paging, 1, limit)
	return
}

// Node type is basic container for anime or manga
type Node struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	MainPicture Picture `json:"main_picture"`
}

type Picture struct {
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
