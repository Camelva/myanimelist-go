package myanimelist

import (
	"testing"
)

func TestMAL_DeleteAnimeFromList(t *testing.T) {
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
			if err := mal.DeleteAnimeFromList(tt.args.animeID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAnimeFromList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMAL_UserAnimeList(t *testing.T) {
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
			got, err := mal.UserAnimeList(tt.args.username, tt.args.status, tt.args.sort, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserAnimeList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Error("UserAnimeList() got nil as result")
				return
			}
		})
	}
}

func TestMAL_UpdateAnimeStatus(t *testing.T) {
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
			got, err := mal.UpdateAnimeStatus(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAnimeStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("UpdateAnimeStatus() got nil as result")
				return
			}

			if got.Status != tt.wantStatus {
				t.Errorf("UpdateAnimeStatus() wrong status - want: %s, got: %s", tt.wantStatus, got.Status)
				return
			}
		})
	}
}
