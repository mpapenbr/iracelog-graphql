//nolint:funlen // ok for test
package storage

import (
	"reflect"
	"testing"

	"github.com/mpapenbr/iracelog-graphql/internal/events"
)

func TestExtractEventSearchKeys(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name    string
		args    args
		want    *events.EventSearchKeys
		wantErr bool
	}{
		{name: "empty", args: args{arg: ""}, wantErr: true},
		{
			name: "team (ci+trim)",
			args: args{arg: "Team: hallo "},
			want: &events.EventSearchKeys{Team: "hallo"},
		},
		{
			name: "car",
			args: args{arg: "cAr: hallo "},
			want: &events.EventSearchKeys{Car: "hallo"},
		},
		{
			name: "driver",
			args: args{arg: "driver:hallo "},
			want: &events.EventSearchKeys{Driver: "hallo"},
		},
		{
			name: "track",
			args: args{arg: "xx:yy  track:hallo "},
			want: &events.EventSearchKeys{Track: "hallo"},
		},
		{
			name: "track+driver",
			args: args{arg: "track:hallo driver:me"},
			want: &events.EventSearchKeys{Track: "hallo", Driver: "me"},
		},
		{
			name: "track (two words)",
			args: args{arg: "track:hallo track"},
			want: &events.EventSearchKeys{Track: "hallo track"},
		},
		{
			name: "regex special",
			args: args{arg: "team:pgz $114 name: (round)"},
			want: &events.EventSearchKeys{Team: "pgz \\\\$114", Name: "(round)"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractEventSearchKeys(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractEventSearchKeys() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractEventSearchKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
