package myanimelist

import (
	"reflect"
	"testing"
)

func TestMAL_AnimeSearch(t *testing.T) {
	type args struct {
		search   string
		settings PagingSettings
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Search 3 anime with word `world`",
			args:    args{"world", PagingSettings{Limit: 3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.AnimeSearch(tt.args.search, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnimeSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("AnimeSearch() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("AnimeSearch() got = %v, want %v", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

var exampleAnimeSearchResult = &AnimeSearchResult{
	Data: nil,
	Paging: Paging{
		Previous: "https://api.myanimelist.net/v2/anime?offset=0&limit=3&q=piece",
		Next:     "https://api.myanimelist.net/v2/anime?offset=6&limit=3&q=piece",
	},
}

func TestAnimeSearchResult_Next(t *testing.T) {
	type args struct {
		obj *AnimeSearchResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeSearchResult},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next(ExampleMAL)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeSearchResult_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeSearchResult_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeSearchResult_Next() got no data!")
				return
			}
		})
	}
}
func TestAnimeSearchResult_Prev(t *testing.T) {
	type args struct {
		obj *AnimeSearchResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeSearchResult},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev(ExampleMAL)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeSearchResult_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeSearchResult_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeSearchResult_Prev() got no data!")
				return
			}
		})
	}
}

func TestMAL_AnimeDetails(t *testing.T) {
	type args struct {
		animeID int
		fields  []string
	}
	tests := []struct {
		name    string
		args    args
		want    *AnimeDetails
		wantErr bool
	}{
		{
			name: "FMA: Brotherhood",
			args: args{5114, []string{FieldTitle}},
			want: &AnimeDetails{
				ID:    5114,
				Title: "Fullmetal Alchemist: Brotherhood",
				MainPicture: Picture{
					Medium: "https://api-cdn.myanimelist.net/images/anime/1223/96541.jpg",
					Large:  "https://api-cdn.myanimelist.net/images/anime/1223/96541l.jpg",
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.AnimeDetails(tt.args.animeID, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnimeDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AnimeDetails() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMAL_AnimeRanking(t *testing.T) {
	type args struct {
		rankingType string
		settings    PagingSettings
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Top 3 favorite anime",
			args:    args{RankFavorite, PagingSettings{Limit: 3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.AnimeRanking(tt.args.rankingType, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnimeRanking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("AnimeRanking() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("AnimeRanking() got = %v, want %v", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

var exampleAnimeRanking = &AnimeRanking{
	Data: nil,
	Paging: Paging{
		Previous: "https://api.myanimelist.net/v2/anime/ranking?offset=0&limit=3",
		Next:     "https://api.myanimelist.net/v2/anime/ranking?offset=6&limit=3",
	},
}

func TestAnimeRanking_Next(t *testing.T) {
	type args struct {
		obj *AnimeRanking
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeRanking},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next(ExampleMAL)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeRanking_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeRanking_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeRanking_Next() got no data!")
				return
			}
		})
	}
}
func TestAnimeRanking_Prev(t *testing.T) {
	type args struct {
		obj *AnimeRanking
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeRanking},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev(ExampleMAL)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeRanking_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeRanking_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeRanking_Prev() got no data!")
				return
			}
		})
	}
}

func TestMAL_SeasonalAnime(t *testing.T) {
	type args struct {
		year     int
		season   string
		sort     string
		settings PagingSettings
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Get 3 anime from 2015 Fall",
			args:    args{year: 2015, season: SeasonFall, sort: SortByScore, settings: PagingSettings{Limit: 3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.SeasonalAnime(tt.args.year, tt.args.season, tt.args.sort, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("SeasonalAnime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("SeasonalAnime() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("SeasonalAnime() got = %v, want %v", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

var exampleSeasonalAnime = &SeasonalAnime{
	Data: nil,
	Paging: Paging{
		Previous: "https://api.myanimelist.net/v2/anime/season/2015/fall?offset=0&limit=3&sort=anime_score",
		Next:     "https://api.myanimelist.net/v2/anime/season/2015/fall?offset=6&limit=3&sort=anime_score",
	},
}

func TestSeasonalAnime_Next(t *testing.T) {
	type args struct {
		obj *SeasonalAnime
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleSeasonalAnime},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next(ExampleMAL)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestSeasonalAnime_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestSeasonalAnime_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestSeasonalAnime_Next() got no data!")
				return
			}
		})
	}
}
func TestSeasonalAnime_Prev(t *testing.T) {
	type args struct {
		obj *SeasonalAnime
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleSeasonalAnime},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev(ExampleMAL)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestSeasonalAnime_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestSeasonalAnime_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestSeasonalAnime_Prev() got no data!")
				return
			}
		})
	}
}

func TestMAL_SuggestedAnime(t *testing.T) {
	type args struct {
		settings PagingSettings
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Get 3 suggestions",
			args:    args{settings: PagingSettings{Limit: 3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.SuggestedAnime(tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuggestedAnime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("SuggestedAnime() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("SuggestedAnime() got = %v, want %v", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

var exampleSuggestedAnime = &SuggestedAnime{
	Data: nil,
	Paging: Paging{
		Previous: "https://api.myanimelist.net/v2/anime/suggestions?offset=0&limit=3",
		Next:     "https://api.myanimelist.net/v2/anime/suggestions?offset=6&limit=3",
	},
}

func TestSuggestedAnime_Next(t *testing.T) {
	type args struct {
		obj *SuggestedAnime
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleSuggestedAnime},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next(ExampleMAL)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestSuggestedAnime_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestSuggestedAnime_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestSuggestedAnime_Next() got no data!")
				return
			}
		})
	}
}
func TestSuggestedAnime_Prev(t *testing.T) {
	type args struct {
		obj *SuggestedAnime
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleSuggestedAnime},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev(ExampleMAL)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestSuggestedAnime_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestSuggestedAnime_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestSuggestedAnime_Prev() got no data!")
				return
			}
		})
	}
}
