package myanimelist

import (
	"reflect"
	"testing"
)

func TestMAL_Anime_Search(t *testing.T) {
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
			got, err := ExampleMAL.Anime.Search(tt.args.search, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Anime_Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("TestMAL_Anime_Search() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("TestMAL_Anime_Search() got = %v, want %v", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

func generateAnimeSearchResult(mal *MAL) *AnimeSearchResult {
	return &AnimeSearchResult{
		parent: &mal.Anime,
		Data:   nil,
		Paging: Paging{
			Previous: "https://api.myanimelist.net/v2/anime?offset=0&limit=3&q=piece",
			Next:     "https://api.myanimelist.net/v2/anime?offset=6&limit=3&q=piece",
		},
	}
}

func TestAnimeSearchResult_Next(t *testing.T) {
	exampleAnimeSearchResult := generateAnimeSearchResult(ExampleMAL)
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
			got, err := tt.args.obj.Next()
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
	exampleAnimeSearchResult := generateAnimeSearchResult(ExampleMAL)
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
			got, err := tt.args.obj.Prev()
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

func TestMAL_Anime_Details(t *testing.T) {
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
			got, err := ExampleMAL.Anime.Details(tt.args.animeID, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Anime_Details() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestMAL_Anime_Details() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMAL_Anime_Top(t *testing.T) {
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
			got, err := ExampleMAL.Anime.Top(tt.args.rankingType, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Anime_Top() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("TestMAL_Anime_Top() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("TestMAL_Anime_Top() got = %v, want %v", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

func generateAnimeTop(mal *MAL) *AnimeTop {
	return &AnimeTop{
		parent: &mal.Anime,
		Data:   nil,
		Paging: Paging{
			Previous: "https://api.myanimelist.net/v2/anime/ranking?offset=0&limit=3",
			Next:     "https://api.myanimelist.net/v2/anime/ranking?offset=6&limit=3",
		},
	}
}

func TestAnimeTop_Next(t *testing.T) {
	exampleAnimeTop := generateAnimeTop(ExampleMAL)
	type args struct {
		obj *AnimeTop
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeTop},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeTop_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeTop_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeTop_Next() got no data!")
				return
			}
		})
	}
}
func TestAnimeTop_Prev(t *testing.T) {
	exampleAnimeTop := generateAnimeTop(ExampleMAL)
	type args struct {
		obj *AnimeTop
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeTop},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeTop_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeTop_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeTop_Prev() got no data!")
				return
			}
		})
	}
}

func TestMAL_Anime_Seasonal(t *testing.T) {
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
			got, err := ExampleMAL.Anime.Seasonal(tt.args.year, tt.args.season, tt.args.sort, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Anime_Seasonal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("TestMAL_Anime_Seasonal() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("TestMAL_Anime_Seasonal() got = %v, want %v", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

func generateAnimeSeasonal(mal *MAL) *AnimeSeasonal {
	return &AnimeSeasonal{
		parent: &mal.Anime,
		Data:   nil,
		Paging: Paging{
			Previous: "https://api.myanimelist.net/v2/anime/season/2015/fall?offset=0&limit=3&sort=anime_score",
			Next:     "https://api.myanimelist.net/v2/anime/season/2015/fall?offset=6&limit=3&sort=anime_score",
		},
	}
}

func TestAnimeSeasonal_Next(t *testing.T) {
	exampleAnimeSeasonal := generateAnimeSeasonal(ExampleMAL)
	type args struct {
		obj *AnimeSeasonal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeSeasonal},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeSeasonal_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeSeasonal_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeSeasonal_Next() got no data!")
				return
			}
		})
	}
}
func TestAnimeSeasonal_Prev(t *testing.T) {
	exampleAnimeSeasonal := generateAnimeSeasonal(ExampleMAL)
	type args struct {
		obj *AnimeSeasonal
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeSeasonal},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeSeasonal_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeSeasonal_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeSeasonal_Prev() got no data!")
				return
			}
		})
	}
}

func TestMAL_Anime_Suggestions(t *testing.T) {
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
			got, err := ExampleMAL.Anime.Suggestions(tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Anime_Suggestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("TestMAL_Anime_Suggestions() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("TestMAL_Anime_Suggestions() got = %v, want %v", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

func generateAnimeSuggestions(mal *MAL) *AnimeSuggestions {
	return &AnimeSuggestions{
		parent: &mal.Anime,
		Data:   nil,
		Paging: Paging{
			Previous: "https://api.myanimelist.net/v2/anime/suggestions?offset=0&limit=3",
			Next:     "https://api.myanimelist.net/v2/anime/suggestions?offset=6&limit=3",
		},
	}
}

func TestAnimeSuggestions_Next(t *testing.T) {
	exampleAnimeSuggestions := generateAnimeSuggestions(ExampleMAL)
	type args struct {
		obj *AnimeSuggestions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeSuggestions},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeSuggestions_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeSuggestions_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeSuggestions_Next() got no data!")
				return
			}
		})
	}
}
func TestAnimeSuggestions_Prev(t *testing.T) {
	exampleAnimeSuggestions := generateAnimeSuggestions(ExampleMAL)
	type args struct {
		obj *AnimeSuggestions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleAnimeSuggestions},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestAnimeSuggestions_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestAnimeSuggestions_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestAnimeSuggestions_Prev() got no data!")
				return
			}
		})
	}
}
