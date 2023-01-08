package events

import (
	"reflect"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mpapenbr/iracelog-graphql/internal"
	tcpg "github.com/mpapenbr/iracelog-graphql/testsupport/tcpostgres"
	"golang.org/x/exp/slices"
)

// Note: we don't check every attribute here. We just pick some to verify the results for specific requests

type checkData struct {
	id        int
	eventName string
	eventKey  string
}

func extractCheckData(dbData []DbEvent) []checkData {
	ret := make([]checkData, len(dbData))
	for i, item := range dbData {
		ret[i] = checkData{id: item.ID, eventName: item.Name}
	}
	return ret
}

func extractAndSortCheckData(dbData []DbEvent) []checkData {
	ret := extractCheckData(dbData)
	slices.SortFunc(ret, func(a, b checkData) bool { return a.id < b.id })
	return ret
}

// we need pointer to ints some DbPageable parameter. Can't do that inside struct as &1
func intHelper(i int) *int {
	return &i
}

// Testdata: (ordered by event name asc)
// 64 - Mid-Ohio 2022-11-20-1617  (no trailing space in event.name, trailing space in event.data->'info'->name)
// 98 - Mid-Ohio 2022-11-20-1617  (trailing space in event.name, trailing space in event.data->'info'->name)
// 50 - Petite LeMans
// 63 - Suzuka 10h
// 48 - Watkins Geln 2022-10-07-2255

func TestGetALl(t *testing.T) {
	pool := tcpg.SetupTestDb()
	type args struct {
		pool     *pgxpool.Pool
		pageable internal.DbPageable
	}
	tests := []struct {
		name    string
		args    args
		want    []checkData
		wantErr bool
	}{
		{name: "2 results, eventName desc", args: args{pool: pool,
			pageable: internal.DbPageable{Sort: []internal.DbSortArg{{Column: "name", Order: "desc"}}, Limit: intHelper(2)}},
			want: []checkData{{id: 48, eventName: "Watkins Glen 2022-10-07-2255"}, {id: 63, eventName: "Suzuka 10h"}}, wantErr: false},
		{name: "2 results, offset 1, eventName asc", args: args{pool: pool,
			pageable: internal.DbPageable{Sort: []internal.DbSortArg{{Column: "name", Order: "asc"}}, Limit: intHelper(2), Offset: intHelper(1)}},
			want: []checkData{{id: 98, eventName: "Mid-Ohio 2022-11-20-1617 "}, {id: 50, eventName: "Petite Lemans"}}, wantErr: false},
		{name: "2 results, offset 1, eventInfo->name asc,id desc", args: args{pool: pool,
			pageable: internal.DbPageable{Sort: []internal.DbSortArg{{Column: "data->'info'->>'name'", Order: "asc"}, {Column: "id", Order: "desc"}}, Limit: intHelper(2), Offset: intHelper(1)}},
			want: []checkData{{id: 64, eventName: "Mid-Ohio 2022-11-20-1617"}, {id: 50, eventName: "Petite Lemans"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetALl(tt.args.pool, tt.args.pageable)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetALl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			check := extractCheckData(got)
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("GetAll() = %v, want %v", check, tt.want)
			}
		})
	}
}

func TestGetByIds(t *testing.T) {
	pool := tcpg.SetupTestDb()
	type args struct {
		pool *pgxpool.Pool
		ids  []int
	}
	tests := []struct {
		name    string
		args    args
		want    []checkData
		wantErr bool
	}{
		{name: "Get events 50,63", args: args{pool: pool, ids: []int{50, 63}}, wantErr: false,
			want: []checkData{{id: 50, eventName: "Petite Lemans"}, {id: 63, eventName: "Suzuka 10h"}}},
		{name: "empty request", args: args{pool: pool, ids: []int{}}, wantErr: false, want: []checkData{}},
		{name: "unknown ids", args: args{pool: pool, ids: []int{999, 3333}}, wantErr: false, want: []checkData{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetByIds(tt.args.pool, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			check := extractAndSortCheckData(got) // ensure ordered array for DeepEqual
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("GetByIds() = %v, want %v", check, tt.want)
			}
		})
	}
}

func TestGetEventsByTrackIds(t *testing.T) {
	pool := tcpg.SetupTestDb()
	type args struct {
		pool     *pgxpool.Pool
		trackIds []int
		pageable internal.DbPageable
	}
	tests := []struct {
		name    string
		args    args
		want    map[int][]checkData
		wantErr bool
	}{
		{name: "Mid-Ohio,Barber", args: args{pool: pool, trackIds: []int{153, 46},
			pageable: internal.DbPageable{Sort: []internal.DbSortArg{{Column: "id", Order: "asc"}}}},

			want: map[int][]checkData{
				153: {{id: 64, eventName: "Mid-Ohio 2022-11-20-1617"}, {id: 98, eventName: "Mid-Ohio 2022-11-20-1617 "}},
				// Note: tracks without event are not present in the map
			},
		},
		{name: "Mid-Ohio, custom sort", args: args{pool: pool, trackIds: []int{153, 46},
			pageable: internal.DbPageable{Sort: []internal.DbSortArg{{Column: "data->'info'->>'name'", Order: "asc"}, {Column: "id", Order: "desc"}}}},

			want: map[int][]checkData{
				153: {{id: 98, eventName: "Mid-Ohio 2022-11-20-1617 "}, {id: 64, eventName: "Mid-Ohio 2022-11-20-1617"}},
				// Note: tracks without event are not present in the map
			},
		},
		{name: "Watkins Glen, Road Atlanta", args: args{pool: pool, trackIds: []int{434, 127}},

			want: map[int][]checkData{
				127: {{id: 50, eventName: "Petite Lemans"}},
				434: {{id: 48, eventName: "Watkins Glen 2022-10-07-2255"}},
			},
		},
		{name: "no results", args: args{pool: pool, trackIds: []int{999, 333}},
			want: map[int][]checkData{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEventsByTrackIds(tt.args.pool, tt.args.trackIds, tt.args.pageable)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsByTrackIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			check := make(map[int][]checkData, len(got))
			for k, v := range got {
				tmp := make([]DbEvent, len(v))
				for i, item := range v {
					tmp[i] = *item
				}
				check[k] = extractCheckData(tmp)
			}
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("GetEventsByTrackIds() = %v, want %v", check, tt.want)
			}
		})
	}
}
