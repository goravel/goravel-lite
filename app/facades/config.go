package facades

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/facades"
)

func Config() config.Config {
	return facades.App().MakeConfig()
}
