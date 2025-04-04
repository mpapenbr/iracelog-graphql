//nolint:dupl // false positive
package storage

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
	"github.com/mpapenbr/iracelog-graphql/internal/tenant"
)

// contains implementations of storage interface that return a model.Tenant items
//
//nolint:whitespace // editor/linter issue
func (db *DBStorage) ResolveTenant(
	ctx context.Context,
	externalID string,
) (ret int, err error) {
	var uuidArg uuid.UUID
	uuidArg, err = uuid.FromString(externalID)
	if err != nil {
		return 0, err
	}
	var tenantRes *models.Tenant
	if tenantRes, err = tenant.FindByExternalID(db.executor, uuidArg); err != nil {
		return 0, err
	}
	return int(tenantRes.ID), nil
}
