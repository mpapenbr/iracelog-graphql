package tenant

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/stephenafamo/bob"

	"github.com/mpapenbr/iracelog-graphql/internal/db/models"
)

func FindByExternalId(exec bob.Executor, externalId uuid.UUID) (*models.Tenant, error) {
	ret, err := models.Tenants.Query(
		models.SelectWhere.Tenants.ExternalID.EQ(externalId),
	).One(context.Background(), exec)

	return ret, err
}
