package tenant

import (
	"context"
	"errors"
)

type (
	ctxKey         struct{}
	TenantProvider func() (int, error)
)

var ErrNoTenantInContext = errors.New("no tenant in context")

func AddToContext(ctx context.Context, provider TenantProvider) context.Context {
	return context.WithValue(ctx, ctxKey{}, provider)
}

func GetFromContext(ctx context.Context) TenantProvider {
	v := ctx.Value(ctxKey{})
	if v == nil {
		return nil
	}
	if ret, ok := v.(TenantProvider); ok {
		return ret
	}
	return nil
}
