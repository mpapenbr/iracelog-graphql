// Code generated by BobGen psql v0.38.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"fmt"
	"io"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dialect"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/bob/expr"
	"github.com/stephenafamo/bob/mods"
	"github.com/stephenafamo/bob/orm"
	"github.com/stephenafamo/bob/types/pgtypes"
)

// CCarDriver is an object representing the database table.
type CCarDriver struct {
	ID          int32  `db:"id,pk" `
	CCarEntryID int32  `db:"c_car_entry_id" `
	DriverID    int32  `db:"driver_id" `
	Name        string `db:"name" `
	Initials    string `db:"initials" `
	AbbrevName  string `db:"abbrev_name" `
	Irating     int32  `db:"irating" `
	LicLevel    int32  `db:"lic_level" `
	LicSubLevel int32  `db:"lic_sub_level" `
	LicString   string `db:"lic_string" `

	R cCarDriverR `db:"-" `
}

// CCarDriverSlice is an alias for a slice of pointers to CCarDriver.
// This should almost always be used instead of []*CCarDriver.
type CCarDriverSlice []*CCarDriver

// CCarDrivers contains methods to work with the c_car_driver table
var CCarDrivers = psql.NewTablex[*CCarDriver, CCarDriverSlice, *CCarDriverSetter]("", "c_car_driver")

// CCarDriversQuery is a query on the c_car_driver table
type CCarDriversQuery = *psql.ViewQuery[*CCarDriver, CCarDriverSlice]

// cCarDriverR is where relationships are stored.
type cCarDriverR struct {
	CCarEntry *CCarEntry // c_car_driver.c_car_driver_car_entry_id_fkey
}

type cCarDriverColumnNames struct {
	ID          string
	CCarEntryID string
	DriverID    string
	Name        string
	Initials    string
	AbbrevName  string
	Irating     string
	LicLevel    string
	LicSubLevel string
	LicString   string
}

var CCarDriverColumns = buildCCarDriverColumns("c_car_driver")

type cCarDriverColumns struct {
	tableAlias  string
	ID          psql.Expression
	CCarEntryID psql.Expression
	DriverID    psql.Expression
	Name        psql.Expression
	Initials    psql.Expression
	AbbrevName  psql.Expression
	Irating     psql.Expression
	LicLevel    psql.Expression
	LicSubLevel psql.Expression
	LicString   psql.Expression
}

func (c cCarDriverColumns) Alias() string {
	return c.tableAlias
}

func (cCarDriverColumns) AliasedAs(alias string) cCarDriverColumns {
	return buildCCarDriverColumns(alias)
}

func buildCCarDriverColumns(alias string) cCarDriverColumns {
	return cCarDriverColumns{
		tableAlias:  alias,
		ID:          psql.Quote(alias, "id"),
		CCarEntryID: psql.Quote(alias, "c_car_entry_id"),
		DriverID:    psql.Quote(alias, "driver_id"),
		Name:        psql.Quote(alias, "name"),
		Initials:    psql.Quote(alias, "initials"),
		AbbrevName:  psql.Quote(alias, "abbrev_name"),
		Irating:     psql.Quote(alias, "irating"),
		LicLevel:    psql.Quote(alias, "lic_level"),
		LicSubLevel: psql.Quote(alias, "lic_sub_level"),
		LicString:   psql.Quote(alias, "lic_string"),
	}
}

type cCarDriverWhere[Q psql.Filterable] struct {
	ID          psql.WhereMod[Q, int32]
	CCarEntryID psql.WhereMod[Q, int32]
	DriverID    psql.WhereMod[Q, int32]
	Name        psql.WhereMod[Q, string]
	Initials    psql.WhereMod[Q, string]
	AbbrevName  psql.WhereMod[Q, string]
	Irating     psql.WhereMod[Q, int32]
	LicLevel    psql.WhereMod[Q, int32]
	LicSubLevel psql.WhereMod[Q, int32]
	LicString   psql.WhereMod[Q, string]
}

func (cCarDriverWhere[Q]) AliasedAs(alias string) cCarDriverWhere[Q] {
	return buildCCarDriverWhere[Q](buildCCarDriverColumns(alias))
}

func buildCCarDriverWhere[Q psql.Filterable](cols cCarDriverColumns) cCarDriverWhere[Q] {
	return cCarDriverWhere[Q]{
		ID:          psql.Where[Q, int32](cols.ID),
		CCarEntryID: psql.Where[Q, int32](cols.CCarEntryID),
		DriverID:    psql.Where[Q, int32](cols.DriverID),
		Name:        psql.Where[Q, string](cols.Name),
		Initials:    psql.Where[Q, string](cols.Initials),
		AbbrevName:  psql.Where[Q, string](cols.AbbrevName),
		Irating:     psql.Where[Q, int32](cols.Irating),
		LicLevel:    psql.Where[Q, int32](cols.LicLevel),
		LicSubLevel: psql.Where[Q, int32](cols.LicSubLevel),
		LicString:   psql.Where[Q, string](cols.LicString),
	}
}

var CCarDriverErrors = &cCarDriverErrors{
	ErrUniqueCCarDriverPkey: &UniqueConstraintError{
		schema:  "",
		table:   "c_car_driver",
		columns: []string{"id"},
		s:       "c_car_driver_pkey",
	},
}

type cCarDriverErrors struct {
	ErrUniqueCCarDriverPkey *UniqueConstraintError
}

// CCarDriverSetter is used for insert/upsert/update operations
// All values are optional, and do not have to be set
// Generated columns are not included
type CCarDriverSetter struct {
	ID          *int32  `db:"id,pk" `
	CCarEntryID *int32  `db:"c_car_entry_id" `
	DriverID    *int32  `db:"driver_id" `
	Name        *string `db:"name" `
	Initials    *string `db:"initials" `
	AbbrevName  *string `db:"abbrev_name" `
	Irating     *int32  `db:"irating" `
	LicLevel    *int32  `db:"lic_level" `
	LicSubLevel *int32  `db:"lic_sub_level" `
	LicString   *string `db:"lic_string" `
}

func (s CCarDriverSetter) SetColumns() []string {
	vals := make([]string, 0, 10)
	if s.ID != nil {
		vals = append(vals, "id")
	}

	if s.CCarEntryID != nil {
		vals = append(vals, "c_car_entry_id")
	}

	if s.DriverID != nil {
		vals = append(vals, "driver_id")
	}

	if s.Name != nil {
		vals = append(vals, "name")
	}

	if s.Initials != nil {
		vals = append(vals, "initials")
	}

	if s.AbbrevName != nil {
		vals = append(vals, "abbrev_name")
	}

	if s.Irating != nil {
		vals = append(vals, "irating")
	}

	if s.LicLevel != nil {
		vals = append(vals, "lic_level")
	}

	if s.LicSubLevel != nil {
		vals = append(vals, "lic_sub_level")
	}

	if s.LicString != nil {
		vals = append(vals, "lic_string")
	}

	return vals
}

func (s CCarDriverSetter) Overwrite(t *CCarDriver) {
	if s.ID != nil {
		t.ID = *s.ID
	}
	if s.CCarEntryID != nil {
		t.CCarEntryID = *s.CCarEntryID
	}
	if s.DriverID != nil {
		t.DriverID = *s.DriverID
	}
	if s.Name != nil {
		t.Name = *s.Name
	}
	if s.Initials != nil {
		t.Initials = *s.Initials
	}
	if s.AbbrevName != nil {
		t.AbbrevName = *s.AbbrevName
	}
	if s.Irating != nil {
		t.Irating = *s.Irating
	}
	if s.LicLevel != nil {
		t.LicLevel = *s.LicLevel
	}
	if s.LicSubLevel != nil {
		t.LicSubLevel = *s.LicSubLevel
	}
	if s.LicString != nil {
		t.LicString = *s.LicString
	}
}

func (s *CCarDriverSetter) Apply(q *dialect.InsertQuery) {
	q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
		return CCarDrivers.BeforeInsertHooks.RunHooks(ctx, exec, s)
	})

	q.AppendValues(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		vals := make([]bob.Expression, 10)
		if s.ID != nil {
			vals[0] = psql.Arg(*s.ID)
		} else {
			vals[0] = psql.Raw("DEFAULT")
		}

		if s.CCarEntryID != nil {
			vals[1] = psql.Arg(*s.CCarEntryID)
		} else {
			vals[1] = psql.Raw("DEFAULT")
		}

		if s.DriverID != nil {
			vals[2] = psql.Arg(*s.DriverID)
		} else {
			vals[2] = psql.Raw("DEFAULT")
		}

		if s.Name != nil {
			vals[3] = psql.Arg(*s.Name)
		} else {
			vals[3] = psql.Raw("DEFAULT")
		}

		if s.Initials != nil {
			vals[4] = psql.Arg(*s.Initials)
		} else {
			vals[4] = psql.Raw("DEFAULT")
		}

		if s.AbbrevName != nil {
			vals[5] = psql.Arg(*s.AbbrevName)
		} else {
			vals[5] = psql.Raw("DEFAULT")
		}

		if s.Irating != nil {
			vals[6] = psql.Arg(*s.Irating)
		} else {
			vals[6] = psql.Raw("DEFAULT")
		}

		if s.LicLevel != nil {
			vals[7] = psql.Arg(*s.LicLevel)
		} else {
			vals[7] = psql.Raw("DEFAULT")
		}

		if s.LicSubLevel != nil {
			vals[8] = psql.Arg(*s.LicSubLevel)
		} else {
			vals[8] = psql.Raw("DEFAULT")
		}

		if s.LicString != nil {
			vals[9] = psql.Arg(*s.LicString)
		} else {
			vals[9] = psql.Raw("DEFAULT")
		}

		return bob.ExpressSlice(ctx, w, d, start, vals, "", ", ", "")
	}))
}

func (s CCarDriverSetter) UpdateMod() bob.Mod[*dialect.UpdateQuery] {
	return um.Set(s.Expressions()...)
}

func (s CCarDriverSetter) Expressions(prefix ...string) []bob.Expression {
	exprs := make([]bob.Expression, 0, 10)

	if s.ID != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "id")...),
			psql.Arg(s.ID),
		}})
	}

	if s.CCarEntryID != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "c_car_entry_id")...),
			psql.Arg(s.CCarEntryID),
		}})
	}

	if s.DriverID != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "driver_id")...),
			psql.Arg(s.DriverID),
		}})
	}

	if s.Name != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "name")...),
			psql.Arg(s.Name),
		}})
	}

	if s.Initials != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "initials")...),
			psql.Arg(s.Initials),
		}})
	}

	if s.AbbrevName != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "abbrev_name")...),
			psql.Arg(s.AbbrevName),
		}})
	}

	if s.Irating != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "irating")...),
			psql.Arg(s.Irating),
		}})
	}

	if s.LicLevel != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "lic_level")...),
			psql.Arg(s.LicLevel),
		}})
	}

	if s.LicSubLevel != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "lic_sub_level")...),
			psql.Arg(s.LicSubLevel),
		}})
	}

	if s.LicString != nil {
		exprs = append(exprs, expr.Join{Sep: " = ", Exprs: []bob.Expression{
			psql.Quote(append(prefix, "lic_string")...),
			psql.Arg(s.LicString),
		}})
	}

	return exprs
}

// FindCCarDriver retrieves a single record by primary key
// If cols is empty Find will return all columns.
func FindCCarDriver(ctx context.Context, exec bob.Executor, IDPK int32, cols ...string) (*CCarDriver, error) {
	if len(cols) == 0 {
		return CCarDrivers.Query(
			SelectWhere.CCarDrivers.ID.EQ(IDPK),
		).One(ctx, exec)
	}

	return CCarDrivers.Query(
		SelectWhere.CCarDrivers.ID.EQ(IDPK),
		sm.Columns(CCarDrivers.Columns().Only(cols...)),
	).One(ctx, exec)
}

// CCarDriverExists checks the presence of a single record by primary key
func CCarDriverExists(ctx context.Context, exec bob.Executor, IDPK int32) (bool, error) {
	return CCarDrivers.Query(
		SelectWhere.CCarDrivers.ID.EQ(IDPK),
	).Exists(ctx, exec)
}

// AfterQueryHook is called after CCarDriver is retrieved from the database
func (o *CCarDriver) AfterQueryHook(ctx context.Context, exec bob.Executor, queryType bob.QueryType) error {
	var err error

	switch queryType {
	case bob.QueryTypeSelect:
		ctx, err = CCarDrivers.AfterSelectHooks.RunHooks(ctx, exec, CCarDriverSlice{o})
	case bob.QueryTypeInsert:
		ctx, err = CCarDrivers.AfterInsertHooks.RunHooks(ctx, exec, CCarDriverSlice{o})
	case bob.QueryTypeUpdate:
		ctx, err = CCarDrivers.AfterUpdateHooks.RunHooks(ctx, exec, CCarDriverSlice{o})
	case bob.QueryTypeDelete:
		ctx, err = CCarDrivers.AfterDeleteHooks.RunHooks(ctx, exec, CCarDriverSlice{o})
	}

	return err
}

// primaryKeyVals returns the primary key values of the CCarDriver
func (o *CCarDriver) primaryKeyVals() bob.Expression {
	return psql.Arg(o.ID)
}

func (o *CCarDriver) pkEQ() dialect.Expression {
	return psql.Quote("c_car_driver", "id").EQ(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		return o.primaryKeyVals().WriteSQL(ctx, w, d, start)
	}))
}

// Update uses an executor to update the CCarDriver
func (o *CCarDriver) Update(ctx context.Context, exec bob.Executor, s *CCarDriverSetter) error {
	v, err := CCarDrivers.Update(s.UpdateMod(), um.Where(o.pkEQ())).One(ctx, exec)
	if err != nil {
		return err
	}

	o.R = v.R
	*o = *v

	return nil
}

// Delete deletes a single CCarDriver record with an executor
func (o *CCarDriver) Delete(ctx context.Context, exec bob.Executor) error {
	_, err := CCarDrivers.Delete(dm.Where(o.pkEQ())).Exec(ctx, exec)
	return err
}

// Reload refreshes the CCarDriver using the executor
func (o *CCarDriver) Reload(ctx context.Context, exec bob.Executor) error {
	o2, err := CCarDrivers.Query(
		SelectWhere.CCarDrivers.ID.EQ(o.ID),
	).One(ctx, exec)
	if err != nil {
		return err
	}
	o2.R = o.R
	*o = *o2

	return nil
}

// AfterQueryHook is called after CCarDriverSlice is retrieved from the database
func (o CCarDriverSlice) AfterQueryHook(ctx context.Context, exec bob.Executor, queryType bob.QueryType) error {
	var err error

	switch queryType {
	case bob.QueryTypeSelect:
		ctx, err = CCarDrivers.AfterSelectHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeInsert:
		ctx, err = CCarDrivers.AfterInsertHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeUpdate:
		ctx, err = CCarDrivers.AfterUpdateHooks.RunHooks(ctx, exec, o)
	case bob.QueryTypeDelete:
		ctx, err = CCarDrivers.AfterDeleteHooks.RunHooks(ctx, exec, o)
	}

	return err
}

func (o CCarDriverSlice) pkIN() dialect.Expression {
	if len(o) == 0 {
		return psql.Raw("NULL")
	}

	return psql.Quote("c_car_driver", "id").In(bob.ExpressionFunc(func(ctx context.Context, w io.Writer, d bob.Dialect, start int) ([]any, error) {
		pkPairs := make([]bob.Expression, len(o))
		for i, row := range o {
			pkPairs[i] = row.primaryKeyVals()
		}
		return bob.ExpressSlice(ctx, w, d, start, pkPairs, "", ", ", "")
	}))
}

// copyMatchingRows finds models in the given slice that have the same primary key
// then it first copies the existing relationships from the old model to the new model
// and then replaces the old model in the slice with the new model
func (o CCarDriverSlice) copyMatchingRows(from ...*CCarDriver) {
	for i, old := range o {
		for _, new := range from {
			if new.ID != old.ID {
				continue
			}
			new.R = old.R
			o[i] = new
			break
		}
	}
}

// UpdateMod modifies an update query with "WHERE primary_key IN (o...)"
func (o CCarDriverSlice) UpdateMod() bob.Mod[*dialect.UpdateQuery] {
	return bob.ModFunc[*dialect.UpdateQuery](func(q *dialect.UpdateQuery) {
		q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
			return CCarDrivers.BeforeUpdateHooks.RunHooks(ctx, exec, o)
		})

		q.AppendLoader(bob.LoaderFunc(func(ctx context.Context, exec bob.Executor, retrieved any) error {
			var err error
			switch retrieved := retrieved.(type) {
			case *CCarDriver:
				o.copyMatchingRows(retrieved)
			case []*CCarDriver:
				o.copyMatchingRows(retrieved...)
			case CCarDriverSlice:
				o.copyMatchingRows(retrieved...)
			default:
				// If the retrieved value is not a CCarDriver or a slice of CCarDriver
				// then run the AfterUpdateHooks on the slice
				_, err = CCarDrivers.AfterUpdateHooks.RunHooks(ctx, exec, o)
			}

			return err
		}))

		q.AppendWhere(o.pkIN())
	})
}

// DeleteMod modifies an delete query with "WHERE primary_key IN (o...)"
func (o CCarDriverSlice) DeleteMod() bob.Mod[*dialect.DeleteQuery] {
	return bob.ModFunc[*dialect.DeleteQuery](func(q *dialect.DeleteQuery) {
		q.AppendHooks(func(ctx context.Context, exec bob.Executor) (context.Context, error) {
			return CCarDrivers.BeforeDeleteHooks.RunHooks(ctx, exec, o)
		})

		q.AppendLoader(bob.LoaderFunc(func(ctx context.Context, exec bob.Executor, retrieved any) error {
			var err error
			switch retrieved := retrieved.(type) {
			case *CCarDriver:
				o.copyMatchingRows(retrieved)
			case []*CCarDriver:
				o.copyMatchingRows(retrieved...)
			case CCarDriverSlice:
				o.copyMatchingRows(retrieved...)
			default:
				// If the retrieved value is not a CCarDriver or a slice of CCarDriver
				// then run the AfterDeleteHooks on the slice
				_, err = CCarDrivers.AfterDeleteHooks.RunHooks(ctx, exec, o)
			}

			return err
		}))

		q.AppendWhere(o.pkIN())
	})
}

func (o CCarDriverSlice) UpdateAll(ctx context.Context, exec bob.Executor, vals CCarDriverSetter) error {
	if len(o) == 0 {
		return nil
	}

	_, err := CCarDrivers.Update(vals.UpdateMod(), o.UpdateMod()).All(ctx, exec)
	return err
}

func (o CCarDriverSlice) DeleteAll(ctx context.Context, exec bob.Executor) error {
	if len(o) == 0 {
		return nil
	}

	_, err := CCarDrivers.Delete(o.DeleteMod()).Exec(ctx, exec)
	return err
}

func (o CCarDriverSlice) ReloadAll(ctx context.Context, exec bob.Executor) error {
	if len(o) == 0 {
		return nil
	}

	o2, err := CCarDrivers.Query(sm.Where(o.pkIN())).All(ctx, exec)
	if err != nil {
		return err
	}

	o.copyMatchingRows(o2...)

	return nil
}

type cCarDriverJoins[Q dialect.Joinable] struct {
	typ       string
	CCarEntry modAs[Q, cCarEntryColumns]
}

func (j cCarDriverJoins[Q]) aliasedAs(alias string) cCarDriverJoins[Q] {
	return buildCCarDriverJoins[Q](buildCCarDriverColumns(alias), j.typ)
}

func buildCCarDriverJoins[Q dialect.Joinable](cols cCarDriverColumns, typ string) cCarDriverJoins[Q] {
	return cCarDriverJoins[Q]{
		typ: typ,
		CCarEntry: modAs[Q, cCarEntryColumns]{
			c: CCarEntryColumns,
			f: func(to cCarEntryColumns) bob.Mod[Q] {
				mods := make(mods.QueryMods[Q], 0, 1)

				{
					mods = append(mods, dialect.Join[Q](typ, CCarEntries.Name().As(to.Alias())).On(
						to.ID.EQ(cols.CCarEntryID),
					))
				}

				return mods
			},
		},
	}
}

// CCarEntry starts a query for related objects on c_car_entry
func (o *CCarDriver) CCarEntry(mods ...bob.Mod[*dialect.SelectQuery]) CCarEntriesQuery {
	return CCarEntries.Query(append(mods,
		sm.Where(CCarEntryColumns.ID.EQ(psql.Arg(o.CCarEntryID))),
	)...)
}

func (os CCarDriverSlice) CCarEntry(mods ...bob.Mod[*dialect.SelectQuery]) CCarEntriesQuery {
	pkCCarEntryID := make(pgtypes.Array[int32], len(os))
	for i, o := range os {
		pkCCarEntryID[i] = o.CCarEntryID
	}
	PKArgExpr := psql.Select(sm.Columns(
		psql.F("unnest", psql.Cast(psql.Arg(pkCCarEntryID), "integer[]")),
	))

	return CCarEntries.Query(append(mods,
		sm.Where(psql.Group(CCarEntryColumns.ID).OP("IN", PKArgExpr)),
	)...)
}

func (o *CCarDriver) Preload(name string, retrieved any) error {
	if o == nil {
		return nil
	}

	switch name {
	case "CCarEntry":
		rel, ok := retrieved.(*CCarEntry)
		if !ok {
			return fmt.Errorf("cCarDriver cannot load %T as %q", retrieved, name)
		}

		o.R.CCarEntry = rel

		if rel != nil {
			rel.R.CCarDrivers = CCarDriverSlice{o}
		}
		return nil
	default:
		return fmt.Errorf("cCarDriver has no relationship %q", name)
	}
}

type cCarDriverPreloader struct {
	CCarEntry func(...psql.PreloadOption) psql.Preloader
}

func buildCCarDriverPreloader() cCarDriverPreloader {
	return cCarDriverPreloader{
		CCarEntry: func(opts ...psql.PreloadOption) psql.Preloader {
			return psql.Preload[*CCarEntry, CCarEntrySlice](orm.Relationship{
				Name: "CCarEntry",
				Sides: []orm.RelSide{
					{
						From: TableNames.CCarDrivers,
						To:   TableNames.CCarEntries,
						FromColumns: []string{
							ColumnNames.CCarDrivers.CCarEntryID,
						},
						ToColumns: []string{
							ColumnNames.CCarEntries.ID,
						},
					},
				},
			}, CCarEntries.Columns().Names(), opts...)
		},
	}
}

type cCarDriverThenLoader[Q orm.Loadable] struct {
	CCarEntry func(...bob.Mod[*dialect.SelectQuery]) orm.Loader[Q]
}

func buildCCarDriverThenLoader[Q orm.Loadable]() cCarDriverThenLoader[Q] {
	type CCarEntryLoadInterface interface {
		LoadCCarEntry(context.Context, bob.Executor, ...bob.Mod[*dialect.SelectQuery]) error
	}

	return cCarDriverThenLoader[Q]{
		CCarEntry: thenLoadBuilder[Q](
			"CCarEntry",
			func(ctx context.Context, exec bob.Executor, retrieved CCarEntryLoadInterface, mods ...bob.Mod[*dialect.SelectQuery]) error {
				return retrieved.LoadCCarEntry(ctx, exec, mods...)
			},
		),
	}
}

// LoadCCarEntry loads the cCarDriver's CCarEntry into the .R struct
func (o *CCarDriver) LoadCCarEntry(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if o == nil {
		return nil
	}

	// Reset the relationship
	o.R.CCarEntry = nil

	related, err := o.CCarEntry(mods...).One(ctx, exec)
	if err != nil {
		return err
	}

	related.R.CCarDrivers = CCarDriverSlice{o}

	o.R.CCarEntry = related
	return nil
}

// LoadCCarEntry loads the cCarDriver's CCarEntry into the .R struct
func (os CCarDriverSlice) LoadCCarEntry(ctx context.Context, exec bob.Executor, mods ...bob.Mod[*dialect.SelectQuery]) error {
	if len(os) == 0 {
		return nil
	}

	cCarEntries, err := os.CCarEntry(mods...).All(ctx, exec)
	if err != nil {
		return err
	}

	for _, o := range os {
		for _, rel := range cCarEntries {
			if o.CCarEntryID != rel.ID {
				continue
			}

			rel.R.CCarDrivers = append(rel.R.CCarDrivers, o)

			o.R.CCarEntry = rel
			break
		}
	}

	return nil
}

func attachCCarDriverCCarEntry0(ctx context.Context, exec bob.Executor, count int, cCarDriver0 *CCarDriver, cCarEntry1 *CCarEntry) (*CCarDriver, error) {
	setter := &CCarDriverSetter{
		CCarEntryID: &cCarEntry1.ID,
	}

	err := cCarDriver0.Update(ctx, exec, setter)
	if err != nil {
		return nil, fmt.Errorf("attachCCarDriverCCarEntry0: %w", err)
	}

	return cCarDriver0, nil
}

func (cCarDriver0 *CCarDriver) InsertCCarEntry(ctx context.Context, exec bob.Executor, related *CCarEntrySetter) error {
	cCarEntry1, err := CCarEntries.Insert(related).One(ctx, exec)
	if err != nil {
		return fmt.Errorf("inserting related objects: %w", err)
	}

	_, err = attachCCarDriverCCarEntry0(ctx, exec, 1, cCarDriver0, cCarEntry1)
	if err != nil {
		return err
	}

	cCarDriver0.R.CCarEntry = cCarEntry1

	cCarEntry1.R.CCarDrivers = append(cCarEntry1.R.CCarDrivers, cCarDriver0)

	return nil
}

func (cCarDriver0 *CCarDriver) AttachCCarEntry(ctx context.Context, exec bob.Executor, cCarEntry1 *CCarEntry) error {
	var err error

	_, err = attachCCarDriverCCarEntry0(ctx, exec, 1, cCarDriver0, cCarEntry1)
	if err != nil {
		return err
	}

	cCarDriver0.R.CCarEntry = cCarEntry1

	cCarEntry1.R.CCarDrivers = append(cCarEntry1.R.CCarDrivers, cCarDriver0)

	return nil
}
