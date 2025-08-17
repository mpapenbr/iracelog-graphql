//nolint:funlen // ok for test
package tracks

import (
	"reflect"
	"slices"
	"testing"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/clause"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stretchr/testify/assert"

	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
	tcpg "github.com/mpapenbr/iracelog-graphql/testsupport/tcpostgres"
)

type checkData struct {
	id        int
	trackName string
}

func extractCheckData(dbData models.TrackSlice) []checkData {
	ret := make([]checkData, len(dbData))
	for i, item := range dbData {
		ret[i] = checkData{id: int(item.ID), trackName: item.ShortName}
	}
	return ret
}

func extractAndSortCheckData(dbData models.TrackSlice) []checkData {
	ret := extractCheckData(dbData)
	slices.SortFunc(ret, func(a, b checkData) int { return a.id - b.id })
	return ret
}

// we need pointer to ints some DBPageable parameter. Can't do that inside struct as &1
func intHelper(i int) *int {
	return &i
}

func TestGetALl(t *testing.T) {
	db := bob.NewDB(tcpg.SetupStdlibDB())

	type sortCol struct {
		col psql.Expression
		dir string
	}
	type args struct {
		pageable internal.DBPageable
		sortCols []sortCol
	}

	// composeOrderBy := func(sortOld []internal.DBSortArg) *clause.OrderBy {

	tests := []struct {
		name    string
		args    args
		want    []checkData
		wantErr bool
	}{
		// Testing only makes sense with predictable results (-> needs sorting).
		// We pick two sort colums with different sortings.

		{
			name: "2 results, displayShort asc",
			args: args{
				pageable: internal.DBPageable{Limit: intHelper(2)},
				sortCols: []sortCol{{models.Tracks.Columns.ShortName, "asc"}},
			},
			want: []checkData{
				{id: 268, trackName: "24 Heures"},
				{id: 345, trackName: "Barcelona"},
			},
			wantErr: false,
		},
		{
			name: "2 results, trackLength desc",
			args: args{
				pageable: internal.DBPageable{Limit: intHelper(2)},
				sortCols: []sortCol{
					{models.Tracks.Columns.TrackLength, "desc"},
					{models.Tracks.Columns.ID, "desc"}, // we have 3 Spa entries...,
				},
			},
			want: []checkData{
				{id: 268, trackName: "24 Heures"},
				{id: 525, trackName: "Spa"},
			},
			wantErr: false,
		},
		{
			name: "2 results, trackLength, default sorting (asc)",
			args: args{
				pageable: internal.DBPageable{Limit: intHelper(2)},
				sortCols: []sortCol{{models.Tracks.Columns.TrackLength, ""}},
			},
			want: []checkData{
				{id: 106, trackName: "Watkins"},
				{id: 233, trackName: "Donington"},
			},
			wantErr: false,
		},
		{
			name: "2 results, numSectors desc",
			args: args{
				pageable: internal.DBPageable{Limit: intHelper(2)},
				sortCols: []sortCol{
					{
						dialect.NewExpression(
							psql.F("jsonb_array_length", models.Tracks.Columns.Sectors)),
						"desc",
					},
					{ // include ID to have defined order (both have 7 sectors)
						models.Tracks.Columns.ID, "asc",
					},
				},
			},
			want: []checkData{
				{id: 95, trackName: "Sebring"},
				{id: 268, trackName: "24 Heures"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := clause.OrderBy{}
			for _, sortArg := range tt.args.sortCols {
				order.AppendOrder(clause.OrderDef{
					Expression: sortArg.col,
					Direction:  sortArg.dir,
				})
			}
			tt.args.pageable.Sort = &order
			got, err := GetAll(db, tt.args.pageable)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetALl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, len(tt.want), len(got), "number of results do not match")
			check := extractAndSortCheckData(got)
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("GetAll() = %v, want %v", check, tt.want)
			}
		})
	}
}

func TestGetByIDs(t *testing.T) {
	db := tcpg.SetupStdlibDB()

	type args struct {
		ids []int
	}
	tests := []struct {
		name    string
		args    args
		want    []checkData
		wantErr bool
	}{
		{
			name:    "Get tracks 18,106",
			args:    args{ids: []int{18, 106}},
			wantErr: false,
			want: []checkData{
				{id: 18, trackName: "Road America"}, {id: 106, trackName: "Watkins"},
			},
		},
		{
			name:    "empty request",
			args:    args{ids: []int{}},
			wantErr: false,
			want:    []checkData{},
		},
		{
			name:    "unknown ids",
			args:    args{ids: []int{999, 3333}},
			wantErr: false,
			want:    []checkData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetByIDs(bob.NewDB(db), tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			check := extractAndSortCheckData(got)
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("GetByIDs() = %v, want %v", check, tt.want)
			}
		})
	}
}
