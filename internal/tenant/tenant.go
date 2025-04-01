package tenant

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/stephenafamo/bob"

	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
)

func FindByExternalID(exec bob.Executor, externalID uuid.UUID) (*models.Tenant, error) {
	ret, err := models.Tenants.Query(
		models.SelectWhere.Tenants.ExternalID.EQ(externalID),
	).One(context.Background(), exec)

	return ret, err
}
