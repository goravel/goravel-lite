package bootstrap

import (
	"github.com/goravel/framework/foundation"

	"goravel/config"
)

func Boot() {
	foundation.Configure().
		WithConfig(config.Boot).
		WithProviders(Providers()).
		Run()
}
