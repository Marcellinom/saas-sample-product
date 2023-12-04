package modules

import (
	"context"

	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth"
)

func RegisterModules(ctx context.Context, i *do.Injector) {
	auth.SetupModule(ctx, i)
}
