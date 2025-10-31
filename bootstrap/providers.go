package bootstrap

import (
	"github.com/goravel/framework/contracts/foundation"

	"goravel/app/providers"
)

func Providers() []foundation.ServiceProvider {
	return []foundation.ServiceProvider{
		&providers.AppServiceProvider{},
	}
}
