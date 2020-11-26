package myanimelist

import (
	"testing"
)

func TestMAL_Forum_Boards(t *testing.T) {
	mal := ExampleMAL
	got, err := mal.Forum.Boards()
	if err != nil {
		t.Errorf("TestMAL_Forum_Boards() error = %v", err)
		return
	}
	if got == nil {
		t.Errorf("TestMAL_Forum_Boards() got nil as result")
		return
	}
}

func TestMAL_Forum_Search(t *testing.T) {
	type args struct {
		searchOpts ForumSearchSettings
		settings   PagingSettings
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Search 5 topics with `Titan`",
			args: args{
				searchOpts: ForumSearchSettings{Keyword: "Titan"},
				settings:   PagingSettings{Limit: 5},
			},
			wantErr: false,
		},
		{
			name: "Search 3 topics started by Lead Administrator `Kineta`",
			args: args{
				searchOpts: ForumSearchSettings{TopicStarter: "Kineta"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mal := ExampleMAL
			got, err := mal.Forum.Search(tt.args.searchOpts, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Forum_Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("TestMAL_Forum_Search() got nil as result")
				return
			}
		})
	}
}

func TestMAL_Forum_Topic(t *testing.T) {
	type args struct {
		topicID  int
		settings PagingSettings
	}
	tests := []struct {
		name      string
		args      args
		wantTitle string
		wantErr   bool
	}{
		{
			name:      "Get 1 post from topic",
			args:      args{topicID: 1849732, settings: PagingSettings{Limit: 1}},
			wantTitle: "MAL's New Public API Release Date!",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mal := ExampleMAL
			got, err := mal.Forum.Topic(tt.args.topicID, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestMAL_Forum_Topic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("TestMAL_Forum_Topic() got nil as result")
				return
			}
			if got.Data.Title != tt.wantTitle {
				t.Errorf("TestMAL_Forum_Topic() expected title: %v, got: %v", tt.wantTitle, got.Data.Title)
				return
			}
		})
	}
}
