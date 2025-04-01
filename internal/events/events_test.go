//nolint:funlen // ok for test
package events

import (
	"reflect"
	"slices"
	"testing"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/clause"
	"github.com/stephenafamo/bob/dialect/psql"

	"github.com/mpapenbr/iracelog-graphql/internal"
	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
	tcpg "github.com/mpapenbr/iracelog-graphql/testsupport/tcpostgres"
)

// Note: we don't check every attribute here.
// We just pick some to verify the results for specific requests

var testTenantID = 1

type checkData struct {
	id        int
	eventName string
}
type sortCol struct {
	col psql.Expression
	dir string
}

func extractCheckData(dbData models.EventSlice) []checkData {
	ret := make([]checkData, len(dbData))
	for i, item := range dbData {
		ret[i] = checkData{id: int(item.ID), eventName: item.Name}
	}
	return ret
}

func extractAndSortCheckData(dbData models.EventSlice) []checkData {
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

	type args struct {
		pageable internal.DBPageable
		sortCols []sortCol
	}
	tests := []struct {
		name    string
		args    args
		want    []checkData
		wantErr bool
	}{
		{
			name: "2 results, eventName desc", args: args{
				pageable: internal.DBPageable{
					SortOld: []internal.DBSortArg{{Column: "name", Order: "desc"}},

					Limit: intHelper(2),
				},
				sortCols: []sortCol{{models.EventColumns.Name, "desc"}},
			},
			want: []checkData{
				{id: 8, eventName: "VRPC Sprint Zandvoort"},
				{id: 9, eventName: "VRPC Main Zandvoort"},
			},
			wantErr: false,
		},
		{
			name: "2 results, offset 1, eventName asc", args: args{
				pageable: internal.DBPageable{
					SortOld: []internal.DBSortArg{{Column: "name", Order: "asc"}},
					Limit:   intHelper(2),
					Offset:  intHelper(1),
				},
				sortCols: []sortCol{{models.EventColumns.Name, "asc"}},
			},
			want: []checkData{
				{id: 4, eventName: "6 Hrs of the Glen"},
				{id: 17, eventName: "Bathurst 12 Hour"},
			},
			wantErr: false,
		},
		{
			name: "2 results, offset 1, name asc,id desc", args: args{
				pageable: internal.DBPageable{
					SortOld: []internal.DBSortArg{
						{Column: "name", Order: "desc"}, {Column: "id", Order: "desc"},
					},
					Limit:  intHelper(2),
					Offset: intHelper(1),
				},
				sortCols: []sortCol{
					{models.EventColumns.Name, "desc"},
					{models.EventColumns.ID, "desc"},
				},
			},
			want: []checkData{
				{id: 9, eventName: "VRPC Main Zandvoort"},
				{id: 14, eventName: "VR e.V. Christmas 500"},
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
			got, err := GetALl(db, testTenantID, tt.args.pageable)
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

func TestGetByIDs(t *testing.T) {
	db := bob.NewDB(tcpg.SetupStdlibDB())
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
			name:    "Get events 4,5",
			args:    args{ids: []int{4, 5}},
			wantErr: false,
			want: []checkData{
				{id: 4, eventName: "6 Hrs of the Glen"},
				{id: 5, eventName: "2024 Spa 24"},
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
			got, err := GetByIDs(db, testTenantID, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			check := extractAndSortCheckData(got) // ensure ordered array for DeepEqual
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("GetByIDs() = %v, want %v", check, tt.want)
			}
		})
	}
}

func TestGetEventsByTrackIDs(t *testing.T) {
	db := bob.NewDB(tcpg.SetupStdlibDB())
	type args struct {
		trackIDs []int
		pageable internal.DBPageable
		sortCols []sortCol
	}
	tests := []struct {
		name    string
		args    args
		want    map[int][]checkData
		wantErr bool
	}{
		{
			name: "Daytona,Zolder", args: args{
				trackIDs: []int{192, 199},
				pageable: internal.DBPageable{
					SortOld: []internal.DBSortArg{{Column: "id", Order: "asc"}},
				},
				sortCols: []sortCol{{models.EventColumns.ID, "asc"}},
			},

			want: map[int][]checkData{
				192: {
					{id: 15, eventName: "Roar Before The 24"},
					{id: 16, eventName: "Daytona 24"},
				},
				// Note: tracks without event are not present in the map
			},
		},
		{
			name: "Daytona, custom sort", args: args{
				trackIDs: []int{192, 199},
				pageable: internal.DBPageable{
					SortOld: []internal.DBSortArg{
						{Column: "name", Order: "asc"},
						{Column: "id", Order: "desc"},
					},
				},
				sortCols: []sortCol{
					{models.EventColumns.Name, "asc"},
					{models.EventColumns.ID, "desc"},
				},
			},

			want: map[int][]checkData{
				192: {
					{id: 16, eventName: "Daytona 24"},
					{id: 15, eventName: "Roar Before The 24"},
				},
				// Note: tracks without event are not present in the map
			},
		},
		{
			name: "Sebring, Road Atlanta",
			args: args{trackIDs: []int{95, 127}},

			want: map[int][]checkData{
				95:  {{id: 13, eventName: "GT Endurance"}},
				127: {{id: 10, eventName: "Petit Le Mans 2024"}},
			},
		},
		{
			name: "no results", args: args{trackIDs: []int{999, 333}},
			want: map[int][]checkData{},
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
			got, err := GetEventsByTrackIDs(
				db,
				testTenantID,
				tt.args.trackIDs,
				tt.args.pageable)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventsByTrackIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			check := make(map[int][]checkData, len(got))
			for k, v := range got {
				tmp := make([]*models.Event, len(v))
				copy(tmp, v)
				check[k] = extractCheckData(tmp)
			}
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("GetEventsByTrackIDs() = %v, want %v", check, tt.want)
			}
		})
	}
}

func TestSimpleSearchEvents(t *testing.T) {
	db := bob.NewDB(tcpg.SetupStdlibDB())
	type args struct {
		searchArg string
	}
	tests := []struct {
		name    string
		args    args
		want    []checkData
		wantErr bool
	}{
		{
			name: "Daytona (Track)", args: args{searchArg: "Daytona"},
			want: []checkData{
				{id: 15, eventName: "Roar Before The 24"},
				{id: 16, eventName: "Daytona 24"},
			},
		},
		{
			name: "Papen (Driver)", args: args{searchArg: "Papen"},
			want: []checkData{
				{id: 14, eventName: "VR e.V. Christmas 500"},
			},
		},
		{
			name: "NSX (Car)", args: args{searchArg: "NSX"},
			want: []checkData{
				{id: 16, eventName: "Daytona 24"},
				{id: 17, eventName: "Bathurst 12 Hour"},
			},
		},
		{
			name: "Alpine (Team)", args: args{searchArg: "Alpine"},
			want: []checkData{
				{id: 5, eventName: "2024 Spa 24"},
				{id: 14, eventName: "VR e.V. Christmas 500"},
			},
		},
	}
	order := clause.OrderBy{}
	order.AppendOrder(clause.OrderDef{
		Expression: models.EventColumns.ID,
		Direction:  "asc",
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SimpleEventSearch(db,
				testTenantID,
				tt.args.searchArg,
				internal.DBPageable{
					Sort: &order,
				})
			if (err != nil) != tt.wantErr {
				t.Errorf("TestSimpleSearchEvents() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			check := extractCheckData(got)
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("TestSimpleSearchEvents() = %v, want %v", check, tt.want)
			}
		})
	}
}

func TestAdvancedEventSearch(t *testing.T) {
	db := bob.NewDB(tcpg.SetupStdlibDB())
	type args struct {
		search EventSearchKeys
	}
	tests := []struct {
		name    string
		args    args
		want    []checkData
		wantErr bool
	}{
		{
			name: "Name", args: args{search: EventSearchKeys{Name: "Petit"}},
			want: []checkData{{id: 10, eventName: "Petit Le Mans 2024"}},
		},
		{
			name: "Track", args: args{search: EventSearchKeys{Track: "Daytona"}},
			want: []checkData{
				{id: 15, eventName: "Roar Before The 24"},
				{id: 16, eventName: "Daytona 24"},
			},
		},
		{
			name: "Track+Name",
			args: args{search: EventSearchKeys{Track: "Daytona", Name: "Roar"}},
			want: []checkData{
				{id: 15, eventName: "Roar Before The 24"},
			},
		},

		{
			name: "Car", args: args{search: EventSearchKeys{Car: "NSX"}},
			want: []checkData{
				{id: 16, eventName: "Daytona 24"},
				{id: 17, eventName: "Bathurst 12 Hour"},
			},
		},

		{
			name: "NonExisting Combo",
			args: args{search: EventSearchKeys{Team: "biela", Car: "Ferrari"}},
			want: []checkData{},
		},
	}
	order := clause.OrderBy{}
	order.AppendOrder(clause.OrderDef{
		Expression: models.EventColumns.ID,
		Direction:  "asc",
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AdvancedEventSearch(db,
				testTenantID,
				&tt.args.search,
				internal.DBPageable{
					Sort: &order,
				})
			if (err != nil) != tt.wantErr {
				t.Errorf("AdvancedEventSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			check := extractCheckData(got)
			if !reflect.DeepEqual(check, tt.want) {
				t.Errorf("AdvancedEventSearch() = %v, want %v", check, tt.want)
			}
		})
	}
}
