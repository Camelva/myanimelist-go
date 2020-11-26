package myanimelist

import (
	"testing"
)

func TestMAL_Anime_List_Remove(t *testing.T) {
	type args struct {
		animeID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Delete FMA:Brotherhood",
			args:    args{5114},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mal := ExampleMAL
			if err := mal.Anime.List.Remove(tt.args.animeID); (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Anime_List_Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMAL_Anime_List_User(t *testing.T) {
	type args struct {
		username string
		status   string
		sort     string
		settings PagingSettings
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Get list with up to 10 planned anime and sort by update date",
			args: args{
				username: "",
				status:   StatusWatching,
				sort:     SortListByUpdateDate,
				settings: PagingSettings{Limit: 10},
			},
			wantErr: false,
		},
		{
			name: "Get 3 anime with no more parameters",
			args: args{
				status:   "",
				sort:     "",
				settings: PagingSettings{Limit: 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mal := ExampleMAL
			got, err := mal.Anime.List.User(tt.args.username, tt.args.status, tt.args.sort, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserAnimeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Error("TestMAL_Anime_List_User() got nil as result")
				return
			}
		})
	}
}

func TestMAL_Anime_List_Update(t *testing.T) {
	type args struct {
		config AnimeConfig
	}
	tests := []struct {
		name       string
		args       args
		wantStatus string
		wantErr    bool
	}{
		{
			name: "Update FMA:Brotherhood status",
			args: args{
				config: NewAnimeConfig(5114).
					SetStatus(StatusCompleted).
					SetScore(10).
					SetTags("some", "random", "tags").
					SetWatchedEpisodes(999).
					SetComment("comment"),
			},
			wantStatus: StatusCompleted,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mal := ExampleMAL
			t.Logf("%+v", tt.args.config)
			got, err := mal.Anime.List.Update(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Anime_List_Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("TestMAL_Anime_List_Update() got nil as result")
				return
			}

			if got.Status != tt.wantStatus {
				t.Errorf("TestMAL_Anime_List_Update() wrong status - want: %s, got: %s", tt.wantStatus, got.Status)
				return
			}
		})
	}
}
