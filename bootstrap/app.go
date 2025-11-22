package bootstrap

import (
	"github.com/goravel/framework/foundation"

	"goravel/config"
)

func Boot() {
	foundation.Setup().
		WithConfig(config.Boot).
		WithProviders(Providers()).
		Run()
}
