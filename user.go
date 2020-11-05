package myanimelist

import (
	"net/http"
	"net/url"
	"time"
)

// UserInformation can only retrieve info about current user for now.
func (mal *MAL) UserInformation() (*UserInformation, error) {
	method := http.MethodGet
	path := "./users/@me"
	data := url.Values{
		"fields": {"anime_statistics"},
	}

	userInfo := new(UserInformation)
	if err := mal.request(userInfo, method, path, data); err != nil {
		return nil, err
	}

	return userInfo, nil
}

type UserInformation struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Location        string          `json:"location"`
	JoinedAt        time.Time       `json:"joined_at"`
	AnimeStatistics AnimeStatistics `json:"anime_statistics"`
}

type AnimeStatistics struct {
	NumItemsWatching    int     `json:"num_items_watching"`
	NumItemsCompleted   int     `json:"num_items_completed"`
	NumItemsOnHold      int     `json:"num_items_on_hold"`
	NumItemsDropped     int     `json:"num_items_dropped"`
	NumItemsPlanToWatch int     `json:"num_items_plan_to_watch"`
	NumItems            int     `json:"num_items"`
	NumDaysWatched      float64 `json:"num_days_watched"`
	NumDaysWatching     float64 `json:"num_days_watching"`
	NumDaysCompleted    float64 `json:"num_days_completed"`
	NumDaysOnHold       float64 `json:"num_days_on_hold"`
	NumDaysDropped      float64 `json:"num_days_dropped"`
	NumDays             float64 `json:"num_days"`
	NumEpisodes         int     `json:"num_episodes"`
	NumTimesRewatched   int     `json:"num_times_rewatched"`
	MeanScore           float64 `json:"mean_score"`
}
