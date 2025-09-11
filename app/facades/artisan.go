package facades

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/facades"
)

func Artisan() console.Artisan {
	return facades.App().MakeArtisan()
}
