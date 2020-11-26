package myanimelist

import (
	"testing"
)

func TestMAL_Manga_List_Remove(t *testing.T) {
	type args struct {
		ID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Delete FMA",
			args:    args{25},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mal := ExampleMAL
			if err := mal.Manga.List.Remove(tt.args.ID); (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Manga_List_Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMAL_Manga_List_User(t *testing.T) {
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
			name: "Get list with up to 10 planned manga and sort by update date",
			args: args{
				username: "",
				status:   StatusReading,
				sort:     SortListByUpdateDate,
				settings: PagingSettings{Limit: 10},
			},
			wantErr: false,
		},
		{
			name: "Get 3 manga with no more parameters",
			args: args{
				username: "",
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
			got, err := mal.Manga.List.User(tt.args.username, tt.args.status, tt.args.sort, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Manga_List_User() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestMAL_Manga_List_User() got nil as result")
				return
			}
		})
	}
}

func TestMAL_Manga_List_Update(t *testing.T) {
	type args struct {
		config MangaConfig
	}
	tests := []struct {
		name       string
		args       args
		wantStatus string
		wantErr    bool
	}{
		{
			name: "Update FMA status",
			args: args{
				config: NewMangaConfig(25).
					SetStatus(StatusCompleted).
					SetScore(10).
					SetTags("some", "random", "tags").
					SetChaptersRead(999).
					SetComment("comment"),
			},
			wantStatus: StatusCompleted,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mal := ExampleMAL
			got, err := mal.Manga.List.Update(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Manga_List_Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("TestMAL_Manga_List_Update() got nil as result")
				return
			}

			if got.Status != tt.wantStatus {
				t.Errorf("TestMAL_Manga_List_Update() wrong status - want: %s, got: %s", tt.wantStatus, got.Status)
				return
			}

			t.Logf("%#v", got)
		})
	}
}
