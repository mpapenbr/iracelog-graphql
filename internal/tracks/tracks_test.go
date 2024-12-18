package tracks

import (
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mpapenbr/iracelog-graphql/internal"
	tcpg "github.com/mpapenbr/iracelog-graphql/testsupport/tcpostgres"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

type checkData struct {
	id        int
	trackName string
}

func extractCheckData(dbData []DbTrack) []checkData {
	ret := make([]checkData, len(dbData))
	for i, item := range dbData {
		ret[i] = checkData{id: item.ID, trackName: item.Data.ShortName}
	}
	return ret
}

func extractAndSortCheckData(dbData []DbTrack) []checkData {
	ret := extractCheckData(dbData)
	slices.SortFunc(ret, func(a, b checkData) bool { return a.id < b.id })
	return ret
}

// we need pointer to ints some DbPageable parameter. Can't do that inside struct as &1
func intHelper(i int) *int {
	return &i
}

func TestGetALl(t *testing.T) {
	pool := tcpg.SetupTestDb()

	// dbStorage := storage.NewDbStorageWithPool(pool)
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
		// Testing only makes sense with predictable results (-> needs sorting). We pick two sort colums with different sortings.

		{
			name: "2 results, displayShort asc", args: args{pool: pool, pageable: internal.DbPageable{Sort: []internal.DbSortArg{{Column: "data->>'trackDisplayShortName'", Order: "asc"}}, Limit: intHelper(2)}},
			want: []checkData{{id: 46, trackName: "Barber"}, {id: 145, trackName: "Brands Hatch"}}, wantErr: false,
		},
		{
			name: "2 results, trackLength desc", args: args{pool: pool, pageable: internal.DbPageable{Sort: []internal.DbSortArg{{Column: "data->'trackLength'", Order: "desc"}}, Limit: intHelper(2)}},
			want: []checkData{{id: 262, trackName: "Gesamtstrecke VLN"}, {id: 341, trackName: "Silverstone"}}, wantErr: false,
		},
		{
			name: "2 results, trackLength, default sorting (asc)", args: args{pool: pool, pageable: internal.DbPageable{Sort: []internal.DbSortArg{{Column: "data->'trackLength'"}}, Limit: intHelper(2)}},
			want: []checkData{{id: 153, trackName: "Mid-Ohio"}, {id: 46, trackName: "Barber"}}, wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetALl(tt.args.pool, tt.args.pageable)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetALl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, len(got), len(tt.want), "number of results do not match")
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
		{name: "Get tracks 46,153", args: args{pool: pool, ids: []int{46, 153}}, wantErr: false, want: []checkData{{id: 46, trackName: "Barber"}, {id: 153, trackName: "Mid-Ohio"}}},
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
			check := extractAndSortCheckData(got)
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("GetByIds() = %v, want %v", check, tt.want)
			}
		})
	}
}
