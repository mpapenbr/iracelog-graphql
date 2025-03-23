// Code generated by BobGen psql v0.30.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"database/sql"
	"database/sql/driver"

	"github.com/gofrs/uuid/v5"
	mytypes "github.com/mpapenbr/iracelog-graphql/internal/db/mytypes"
	"github.com/shopspring/decimal"
	"github.com/stephenafamo/bob"
)

// Make sure the type CCar runs hooks after queries
var _ bob.HookableType = &CCar{}

// Make sure the type CCarDriver runs hooks after queries
var _ bob.HookableType = &CCarDriver{}

// Make sure the type CCarEntry runs hooks after queries
var _ bob.HookableType = &CCarEntry{}

// Make sure the type CCarTeam runs hooks after queries
var _ bob.HookableType = &CCarTeam{}

// Make sure the type Event runs hooks after queries
var _ bob.HookableType = &Event{}

// Make sure the type Tenant runs hooks after queries
var _ bob.HookableType = &Tenant{}

// Make sure the type Track runs hooks after queries
var _ bob.HookableType = &Track{}

// Make sure the type decimal.Decimal satisfies database/sql.Scanner
var _ sql.Scanner = (*decimal.Decimal)(nil)

// Make sure the type decimal.Decimal satisfies database/sql/driver.Valuer
var _ driver.Valuer = *new(decimal.Decimal)

// Make sure the type mytypes.EventSessionSlice satisfies database/sql.Scanner
var _ sql.Scanner = (*mytypes.EventSessionSlice)(nil)

// Make sure the type mytypes.EventSessionSlice satisfies database/sql/driver.Valuer
var _ driver.Valuer = *new(mytypes.EventSessionSlice)

// Make sure the type uuid.UUID satisfies database/sql.Scanner
var _ sql.Scanner = (*uuid.UUID)(nil)

// Make sure the type uuid.UUID satisfies database/sql/driver.Valuer
var _ driver.Valuer = *new(uuid.UUID)

// Make sure the type mytypes.SectorSlice satisfies database/sql.Scanner
var _ sql.Scanner = (*mytypes.SectorSlice)(nil)

// Make sure the type mytypes.SectorSlice satisfies database/sql/driver.Valuer
var _ driver.Valuer = *new(mytypes.SectorSlice)
