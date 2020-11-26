package myanimelist

import (
	"reflect"
	"testing"
)

func TestMAL_Manga_Search(t *testing.T) {
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
			name:    "Search 3 manga with word `world`",
			args:    args{search: "world", settings: PagingSettings{Limit: 3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.Manga.Search(tt.args.search, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Manga_Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("TestMAL_Manga_Search() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("TestMAL_Manga_Search() got = %d, want %d", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

func generateExampleMangaSearchResult(mal *MAL) *MangaSearchResult {
	return &MangaSearchResult{
		parent: &mal.Manga,
		Data:   nil,
		Paging: Paging{
			Previous: "https://api.myanimelist.net/v2/manga?offset=0&limit=3&q=piece",
			Next:     "https://api.myanimelist.net/v2/manga?offset=6&limit=3&q=piece",
		},
	}
}

func TestMangaSearchResult_Next(t *testing.T) {
	exampleMangaSearchResult := generateExampleMangaSearchResult(ExampleMAL)
	type args struct {
		obj *MangaSearchResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleMangaSearchResult},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMangaSearchResult_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestMangaSearchResult_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestMangaSearchResult_Next() got no data!")
				return
			}
		})
	}
}
func TestMangaSearchResult_Prev(t *testing.T) {
	exampleMangaSearchResult := generateExampleMangaSearchResult(ExampleMAL)
	type args struct {
		obj *MangaSearchResult
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleMangaSearchResult},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMangaSearchResult_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestMangaSearchResult_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestMangaSearchResult_Prev() got no data!")
				return
			}
		})
	}
}

func TestMAL_Manga_Details(t *testing.T) {
	type args struct {
		mangaID int
		fields  []string
	}
	tests := []struct {
		name    string
		args    args
		want    *MangaDetails
		wantErr bool
	}{
		{
			name: "Berserk",
			args: args{mangaID: 2, fields: []string{FieldTitle}},
			want: &MangaDetails{
				ID:    2,
				Title: "Berserk",
				MainPicture: Picture{
					Medium: "https://api-cdn.myanimelist.net/images/manga/1/157931.jpg",
					Large:  "https://api-cdn.myanimelist.net/images/manga/1/157931l.jpg",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.Manga.Details(tt.args.mangaID, tt.args.fields...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Manga_Details() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TestMAL_Manga_Details() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMAL_Manga_Top(t *testing.T) {
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
			name:    "Top 3 favorite manga",
			args:    args{RankFavorite, PagingSettings{Limit: 3}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExampleMAL.Manga.Top(tt.args.rankingType, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Manga_Top() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("TestMAL_Manga_Top() got nil as result")
				return
			}
			if len(got.Data) < tt.args.settings.Limit {
				t.Errorf("TestMAL_Manga_Top() got = %d, want %d", len(got.Data), tt.args.settings.Limit)
			}
		})
	}
}

func generateExampleMangaTop(mal *MAL) *MangaTop {
	return &MangaTop{
		parent: &mal.Manga,
		Data:   nil,
		Paging: Paging{
			Previous: "https://api.myanimelist.net/v2/manga/ranking?offset=0&limit=3",
			Next:     "https://api.myanimelist.net/v2/manga/ranking?offset=6&limit=3",
		},
	}
}

func TestMangaTop_Next(t *testing.T) {
	exampleMangaTop := generateExampleMangaTop(ExampleMAL)
	type args struct {
		obj *MangaTop
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleMangaTop},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Next()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMangaTop_Next() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestMangaTop_Next() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestMangaTop_Next() got no data!")
				return
			}
		})
	}
}
func TestMangaTop_Prev(t *testing.T) {
	exampleMangaTop := generateExampleMangaTop(ExampleMAL)
	type args struct {
		obj *MangaTop
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Basic example",
			args:    args{exampleMangaTop},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.obj.Prev()
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMangaTop_Prev() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestMangaTop_Prev() got nil")
				return
			}
			if len(got.Data) < 1 {
				t.Error("TestMangaTop_Prev() got no data!")
				return
			}
		})
	}
}
