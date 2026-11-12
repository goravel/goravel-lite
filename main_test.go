package main

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/path"
	"github.com/stretchr/testify/suite"

	"goravel/app/facades"
	"goravel/bootstrap"
)

// MainTestSuite installs packages on the framework's release line, the one go.mod
// pins, and never with --dev: --dev fetches drivers from master, and master
// drivers call framework APIs that line does not have yet. The install helpers
// below hold that rule, so tests only name what they install.
type MainTestSuite struct {
	suite.Suite
	baseline   *workspaceBaseline
	restoreErr error
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (s *MainTestSuite) SetupSuite() {
	bootstrap.Boot()

	// There is no .env file in the testing environment, so we need to initialize the configuration first.
	facades.Config().Add("app.name", "Goravel")
	facades.Config().Add("app.env", "local")
	facades.Config().Add("app.debug", "true")

	baseline, err := newWorkspaceBaseline(mutableTestPaths)
	s.Require().NoError(err)

	s.baseline = baseline
}

func (s *MainTestSuite) SetupTest() {
	if s.restoreErr != nil {
		s.T().Skipf("workspace restore failed in a previous test: %v", s.restoreErr)
	}
}

func (s *MainTestSuite) TearDownTest() {
	if s.baseline == nil || s.restoreErr != nil {
		return
	}

	if err := s.baseline.Restore(); err != nil {
		s.restoreErr = err
		s.T().Fatalf("restore workspace: %v", err)
	}
}

func (s *MainTestSuite) TearDownSuite() {
	if s.baseline == nil {
		return
	}

	defer func() {
		s.baseline = nil
	}()

	s.Require().NoError(s.baseline.Close())
}

func (s *MainTestSuite) TestPackageInstall_All() {
	s.installAllFacades()

	s.False(facades.Process().Run("go", "run", ".", "artisan").Failed())

	s.uninstallPackages(
		"AI",
		"Auth",
		"Testing",
		"Grpc",
		"Hash",
		"Route",
		"Http",
		"View",
		"Session",
		"Storage",
		"Validation",
		"Lang",
		"Mail",
		"Crypt",
		"RateLimiter",
		"Schedule",
		"Gate",
		"Event",
		"Queue",
		"Schema",
		"Seeder",
		"DB",
		"Orm",
		"Cache",
		"Telemetry",
		"Log",
	)
}

func (s *MainTestSuite) TestPackageInstall_AI() {
	s.installFacade("AI")
	s.FileExists(path.Facade("ai.go"))
	s.FileExists(path.Config("ai.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&ai.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall AI"))
	s.NoFileExists(path.Facade("ai.go"))
	s.NoFileExists(path.Config("ai.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&ai.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Auth() {
	s.installFacade("Auth")
	s.FileExists(path.Facade("auth.go"))
	s.FileExists(path.Config("auth.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&auth.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Auth"))
	s.NoFileExists(path.Facade("auth.go"))
	s.NoFileExists(path.Config("auth.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&auth.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Cache() {
	s.installFacade("Cache")
	s.FileExists(path.Facade("cache.go"))
	s.FileExists(path.Config("cache.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&cache.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Cache"))
	s.NoFileExists(path.Facade("cache.go"))
	s.NoFileExists(path.Config("cache.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&cache.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Crypt() {
	s.installFacade("Crypt")
	s.FileExists(path.Facade("crypt.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&crypt.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Crypt"))
	s.NoFileExists(path.Facade("crypt.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&crypt.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_DB() {
	s.installFacade("DB")
	s.FileExists(path.Facade("db.go"))
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&postgres.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "postgres"))

	s.NoError(facades.Artisan().Call("package:uninstall DB"))
	s.NoFileExists(path.Facade("db.go"))

	// The Schema facade still exists, so database.go and ServiceProvider should still exist.
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&postgres.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "postgres"))
}

func (s *MainTestSuite) TestPackageInstall_Event() {
	s.installFacade("Event")
	s.FileExists(path.Facade("event.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&event.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Event"))
	s.NoFileExists(path.Facade("event.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&event.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Gate() {
	s.installFacade("Gate")
	s.FileExists(path.Facade("gate.go"))
	s.FileExists(path.Config("auth.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&auth.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Gate"))
	s.NoFileExists(path.Facade("gate.go"))
	s.NoFileExists(path.Config("auth.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&auth.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Grpc() {
	s.installFacade("Grpc")
	s.FileExists(path.Facade("grpc.go"))
	s.FileExists(path.Config("grpc.go"))
	s.FileExists(path.Route("grpc.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&grpc.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("app.go"), ".Grpc"))
	s.True(file.Contains(path.Base(".env.example"), `
GRPC_HOST=
GRPC_PORT=
`))

	s.NoError(facades.Artisan().Call("package:uninstall Grpc"))
	s.NoFileExists(path.Facade("grpc.go"))
	s.NoFileExists(path.Config("grpc.go"))
	s.NoFileExists(path.Route("grpc.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&grpc.ServiceProvider{},"))
	s.False(file.Contains(path.Bootstrap("app.go"), ".Grpc"))
}

func (s *MainTestSuite) TestPackageInstall_Hash() {
	s.installFacade("Hash")
	s.FileExists(path.Facade("hash.go"))
	s.FileExists(path.Config("hashing.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&hash.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Hash"))
	s.NoFileExists(path.Facade("hash.go"))
	s.NoFileExists(path.Config("hashing.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&hash.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Http() {
	s.installFacade("Http")
	s.FileExists(path.Facade("http.go"))
	s.FileExists(path.Config("http.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&http.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Http"))
	s.NoFileExists(path.Facade("http.go"))
	s.NoFileExists(path.Config("http.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&http.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Lang() {
	s.installFacade("Lang")
	s.FileExists(path.Facade("lang.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&translation.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Lang"))
	s.NoFileExists(path.Facade("lang.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&translation.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Log() {
	s.installFacade("Log")
	s.FileExists(path.Facade("log.go"))
	s.FileExists(path.Config("logging.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&log.ServiceProvider{},"))
	s.True(file.Contains(path.Base(".env.example"), `
LOG_CHANNEL=stack
LOG_LEVEL=debug
`))

	s.NoError(facades.Artisan().Call("package:uninstall Log"))
	s.NoFileExists(path.Facade("log.go"))
	s.NoFileExists(path.Config("logging.go"))
	s.NoFileExists(path.Route("logging.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&log.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Mail() {
	s.installFacade("Mail")
	s.FileExists(path.Facade("mail.go"))
	s.FileExists(path.Config("mail.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&mail.ServiceProvider{},"))
	s.True(file.Contains(path.Base(".env.example"), `
MAIL_HOST=
MAIL_PORT=
MAIL_USERNAME=
MAIL_PASSWORD=
MAIL_FROM_ADDRESS=
MAIL_FROM_NAME=
`))

	s.NoError(facades.Artisan().Call("package:uninstall Mail"))
	s.NoFileExists(path.Facade("mail.go"))
	s.NoFileExists(path.Config("mail.go"))
	s.NoFileExists(path.Route("mail.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&mail.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Orm() {
	s.installFacade("Orm")
	s.FileExists(path.Facade("orm.go"))
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))

	// Schema should be uninstalled if Orm wants to be uninstalled.
	s.NoError(facades.Artisan().Call("package:uninstall Schema"))
	s.NoError(facades.Artisan().Call("package:uninstall Orm"))
	s.NoFileExists(path.Facade("orm.go"))
	s.NoFileExists(path.Config("database.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Queue() {
	s.installFacade("Queue")
	s.FileExists(path.Facade("queue.go"))
	s.FileExists(path.Config("queue.go"))
	s.FileExists(path.Migration("20210101000001_create_jobs_table.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&queue.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("app.go"), "WithMigrations(Migrations)"))
	s.True(file.Contains(path.Bootstrap("migrations.go"), "M20210101000001CreateJobsTable{}"))

	s.NoError(facades.Artisan().Call("package:uninstall Queue"))
	s.NoFileExists(path.Facade("queue.go"))
	s.NoFileExists(path.Config("queue.go"))
	s.NoFileExists(path.Migration("20210101000001_create_jobs_table.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&queue.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("app.go"), "WithMigrations(Migrations)"))
	s.False(file.Contains(path.Bootstrap("migrations.go"), "M20210101000001CreateJobsTable{}"))
}

func (s *MainTestSuite) TestPackageInstall_RateLimiter() {
	s.installFacade("RateLimiter")
	s.FileExists(path.Facade("rate_limiter.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&http.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall RateLimiter"))
	s.NoFileExists(path.Facade("rate_limiter.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&http.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Route() {
	s.installFacade("Route")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&route.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&gin.ServiceProvider{},"))
	s.True(file.Contains(path.Bootstrap("app.go"), ".Web"))
	s.True(file.Contains(path.Config("http.go"), "gin"))
	s.FileExists(path.Resource("views", "welcome.tmpl"))
	s.FileExists(path.App("http", "controllers", "user_controller.go"))
	s.FileExists(path.Public(".gitignore"))
	s.FileExists(path.Route("web.go"))
	s.FileExists(path.Facade("route.go"))
	s.FileExists(path.Config("jwt.go"))
	s.FileExists(path.Config("cors.go"))
	s.True(file.Contains(path.Base(".env.example"), `
APP_URL=http://localhost
APP_HOST=127.0.0.1
APP_PORT=3000

JWT_SECRET=
`))

	s.NoError(facades.Artisan().Call("package:uninstall Route"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&route.ServiceProvider{},"))
	s.False(file.Contains(path.Bootstrap("app.go"), ".Web"))
	// The Http facade still exists, so "gin" related content in http.go should still exist.
	s.True(file.Contains(path.Bootstrap("providers.go"), "&gin.ServiceProvider{},"))
	s.True(file.Contains(path.Config("http.go"), "gin"))
	s.NoFileExists(path.Resource())
	s.NoFileExists(path.App("http"))
	s.NoFileExists(path.Public())
	s.NoFileExists(path.Route("web.go"))
	s.NoFileExists(path.Facade("route.go"))
	s.NoFileExists(path.Config("jwt.go"))
	s.NoFileExists(path.Config("cors.go"))
}

func (s *MainTestSuite) TestPackageInstall_Schedule() {
	s.installFacade("Schedule")
	s.FileExists(path.Facade("schedule.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&schedule.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Schedule"))
	s.NoFileExists(path.Facade("schedule.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&schedule.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Schema() {
	s.installFacade("Schema")
	s.FileExists(path.Facade("schema.go"))
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))

	s.NoError(facades.Artisan().Call("package:uninstall Schema"))
	s.NoFileExists(path.Facade("schema.go"))

	// The Orm facade still exists, so database.go and ServiceProvider should still exist.
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
}

func (s *MainTestSuite) TestPackageInstall_Seeder() {
	s.installFacade("Seeder")
	s.FileExists(path.Facade("seeder.go"))
	s.FileExists(path.Config("database.go"))
	s.True(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "default"))
	s.True(file.Contains(path.Config("database.go"), "connections"))
	s.True(file.Contains(path.Config("database.go"), "pool"))
	s.True(file.Contains(path.Config("database.go"), "slow_threshold"))
	s.True(file.Contains(path.Config("database.go"), "migrations"))
	s.True(file.Contains(path.Base(".env.example"), `
DB_HOST=
DB_PORT=
DB_DATABASE=
DB_USERNAME=
DB_PASSWORD=
`))

	s.NoError(facades.Artisan().Call("package:uninstall Seeder"))
	s.NoFileExists(path.Facade("seeder.go"))
	s.NoFileExists(path.Config("database.go"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&database.ServiceProvider{},"))
	s.False(file.Contains(path.Config("database.go"), "default"))
	s.False(file.Contains(path.Config("database.go"), "connections"))
	s.False(file.Contains(path.Config("database.go"), "pool"))
	s.False(file.Contains(path.Config("database.go"), "slow_threshold"))
	s.False(file.Contains(path.Config("database.go"), "migrations"))
}

func (s *MainTestSuite) TestPackageInstall_Session() {
	s.installFacade("Session")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&session.ServiceProvider{},"))
	s.FileExists(path.Config("session.go"))
	s.FileExists(path.Facade("session.go"))
	s.True(file.Contains(path.Base(".env.example"), `
SESSION_DRIVER=file
SESSION_LIFETIME=120
`))

	s.NoError(facades.Artisan().Call("package:uninstall Session"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&session.ServiceProvider{},"))
	s.NoFileExists(path.Config("session.go"))
	s.NoFileExists(path.Facade("session.go"))
}

func (s *MainTestSuite) TestPackageInstall_Storage() {
	s.installFacade("Storage")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&filesystem.ServiceProvider{},"))
	s.FileExists(path.Config("filesystems.go"))
	s.FileExists(path.Facade("storage.go"))

	s.NoError(facades.Artisan().Call("package:uninstall Storage"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&filesystem.ServiceProvider{},"))
	s.NoFileExists(path.Config("filesystems.go"))
	s.NoFileExists(path.Facade("storage.go"))
}

func (s *MainTestSuite) TestPackageInstall_Telemetry() {
	s.installFacade("Telemetry")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&telemetry.ServiceProvider{},"))
	s.FileExists(path.Config("telemetry.go"))
	s.FileExists(path.Facade("telemetry.go"))
	s.True(file.Contains(path.Config("logging.go"), "instrument_name"))

	s.NoError(facades.Artisan().Call("package:uninstall Telemetry"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&telemetry.ServiceProvider{},"))
	s.NoFileExists(path.Config("telemetry.go"))
	s.NoFileExists(path.Facade("telemetry.go"))
	s.False(file.Contains(path.Config("logging.go"), "instrument_name"))
}

func (s *MainTestSuite) TestPackageInstall_Testing() {
	s.installFacade("Testing")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&testing.ServiceProvider{},"))
	s.FileExists(path.Test("test_case.go"))
	s.FileExists(path.Test("feature", "example_test.go"))
	s.FileExists(path.Facade("testing.go"))

	s.NoError(facades.Artisan().Call("package:uninstall Testing"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&testing.ServiceProvider{},"))
	s.NoFileExists(path.Test("test_case.go"))
	s.NoFileExists(path.Test("feature", "example_test.go"))
	s.NoFileExists(path.Facade("testing.go"))
}

func (s *MainTestSuite) TestPackageInstall_Validation() {
	s.installFacade("Validation")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&validation.ServiceProvider{},"))
	s.FileExists(path.Facade("validation.go"))

	s.NoError(facades.Artisan().Call("package:uninstall Validation"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&validation.ServiceProvider{},"))
	s.NoFileExists(path.Facade("validation.go"))
}

func (s *MainTestSuite) TestPackageInstall_View() {
	s.installFacade("View")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&view.ServiceProvider{},"))
	s.FileExists(path.Facade("view.go"))

	s.NoError(facades.Artisan().Call("package:uninstall View"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&view.ServiceProvider{},"))
	s.NoFileExists(path.Facade("view.go"))
}

func (s *MainTestSuite) TestPackageInstall_DBDrivers() {
	s.installFacade("DB")

	s.installDriver("github.com/goravel/mysql")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&mysql.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "mysql"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/mysql"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&mysql.ServiceProvider{},"))
	s.False(file.Contains(path.Config("database.go"), "mysql"))

	s.installDriver("github.com/goravel/sqlserver")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&sqlserver.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "sqlserver"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/sqlserver"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&sqlserver.ServiceProvider{},"))
	s.False(file.Contains(path.Config("database.go"), "sqlserver"))

	s.installDriver("github.com/goravel/sqlite")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&sqlite.ServiceProvider{},"))
	s.True(file.Contains(path.Config("database.go"), "sqlite"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/sqlite"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&sqlite.ServiceProvider{},"))
	s.False(file.Contains(path.Config("database.go"), "sqlite"))
}

func (s *MainTestSuite) TestPackageInstall_FilesystemDrivers() {
	s.installFacade("Storage")

	s.installDriver("github.com/goravel/s3")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&s3.ServiceProvider{},"))
	s.True(file.Contains(path.Config("filesystems.go"), "s3"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/s3"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&s3.ServiceProvider{},"))
	s.False(file.Contains(path.Config("filesystems.go"), "s3"))

	s.installDriver("github.com/goravel/oss")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&oss.ServiceProvider{},"))
	s.True(file.Contains(path.Config("filesystems.go"), "oss"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/oss"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&oss.ServiceProvider{},"))
	s.False(file.Contains(path.Config("filesystems.go"), "oss"))

	s.installDriver("github.com/goravel/cos")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&cos.ServiceProvider{},"))
	s.True(file.Contains(path.Config("filesystems.go"), "cos"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/cos"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&cos.ServiceProvider{},"))
	s.False(file.Contains(path.Config("filesystems.go"), "cos"))

	s.installDriver("github.com/goravel/minio")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&minio.ServiceProvider{},"))
	s.True(file.Contains(path.Config("filesystems.go"), "minio"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/minio"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&minio.ServiceProvider{},"))
	s.False(file.Contains(path.Config("filesystems.go"), "minio"))
}

func (s *MainTestSuite) TestPackageInstall_CacheDrivers() {
	s.installFacade("Cache")
	s.installFacade("Session")
	s.installFacade("Queue")

	s.installDriver("github.com/goravel/redis")
	s.True(file.Contains(path.Bootstrap("providers.go"), "&redis.ServiceProvider{},"))
	s.True(file.Contains(path.Config("cache.go"), "redis"))
	s.True(file.Contains(path.Config("database.go"), "redis"))
	s.True(file.Contains(path.Config("queue.go"), "redis"))
	s.True(file.Contains(path.Config("session.go"), "redis"))

	s.NoError(facades.Artisan().Call("package:uninstall github.com/goravel/redis"))
	s.False(file.Contains(path.Bootstrap("providers.go"), "&redis.ServiceProvider{},"))
	s.False(file.Contains(path.Config("cache.go"), "redis"))
	s.False(file.Contains(path.Config("database.go"), "redis"))
	s.False(file.Contains(path.Config("queue.go"), "redis"))
	s.False(file.Contains(path.Config("session.go"), "redis"))
}

func (s *MainTestSuite) artisan(command string) {
	s.NoError(facades.Artisan().Call(command))
}

func (s *MainTestSuite) installFacade(facade string) {
	s.artisan("package:install " + facade + " --default")
}

func (s *MainTestSuite) installAllFacades() {
	s.artisan("package:install --all --default")
}

// installDriver pins pkg to the framework's release line, because package:install
// only resolves that line for the drivers a facade declares, not for a named
// package. The branch ref, v1.18.x for a v1.18 line, covers drivers with no
// release tag on the line yet.
//
// The ref is picked with go list because package:install prints a failed install
// and still returns nil, so its own error tells us nothing.
func (s *MainTestSuite) installDriver(pkg string) {
	line := releaseLine(support.Version)
	s.Require().NotEmpty(line, "framework version %q has no release line", support.Version)

	var tried []string
	for _, ref := range []string{line, line + ".x"} {
		tried = append(tried, ref)
		if !facades.Process().Run("go", "list", "-m", pkg+"@"+ref).Failed() {
			s.artisan("package:install " + pkg + "@" + ref)
			return
		}
	}

	s.FailNow("install " + pkg + " on the framework release line: no ref resolved, tried " + strings.Join(tried, ", "))
}

// releaseLine turns a framework version into its minor line, v1.18.3 into v1.18.
func releaseLine(version string) string {
	major, rest, ok := strings.Cut(strings.TrimPrefix(version, "v"), ".")
	if !ok {
		return ""
	}

	minor, _, _ := strings.Cut(rest, ".")
	if _, err := strconv.Atoi(major); err != nil {
		return ""
	}
	if _, err := strconv.Atoi(minor); err != nil {
		return ""
	}

	return "v" + major + "." + minor
}

func (s *MainTestSuite) uninstallPackages(packages ...string) {
	for _, pkg := range packages {
		s.artisan("package:uninstall " + pkg)
	}
}

var mutableTestPaths = []string{
	".env.example",
	"go.mod",
	"go.sum",
	"bootstrap",
	"config",
	"app/facades",
	"app/http",
	"database",
	"public",
	"resources",
	"routes",
	"tests",
}

type workspaceBaseline struct {
	roots   []string
	tempDir string
}

func newWorkspaceBaseline(paths []string) (*workspaceBaseline, error) {
	tempDir, err := os.MkdirTemp("", "goravel-lite-workspace-*")
	if err != nil {
		return nil, err
	}

	baseline := &workspaceBaseline{
		roots:   make([]string, 0, len(paths)),
		tempDir: tempDir,
	}

	for _, path := range paths {
		path = filepath.Clean(path)
		baseline.roots = append(baseline.roots, path)

		if err := baseline.capture(path); err != nil {
			_ = baseline.Close()

			return nil, err
		}
	}

	return baseline, nil
}

func (b *workspaceBaseline) capture(path string) error {
	if _, err := os.Lstat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return err
	}

	return copyPath(path, b.snapshotPath(path))
}

func (b *workspaceBaseline) Restore() error {
	for _, root := range b.roots {
		if err := os.RemoveAll(root); err != nil {
			return err
		}
	}

	for _, root := range b.roots {
		snapshotPath := b.snapshotPath(root)
		if _, err := os.Lstat(snapshotPath); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			return err
		}

		if err := copyPath(snapshotPath, root); err != nil {
			return err
		}
	}

	return nil
}

func (b *workspaceBaseline) Close() error {
	return os.RemoveAll(b.tempDir)
}

func (b *workspaceBaseline) snapshotPath(path string) string {
	return filepath.Join(b.tempDir, path)
}

func copyPath(src, dst string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}

	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		if err := mkdirParent(dst); err != nil {
			return err
		}

		link, err := os.Readlink(src)
		if err != nil {
			return err
		}

		return os.Symlink(link, dst)
	case info.IsDir():
		if err := os.MkdirAll(dst, info.Mode().Perm()); err != nil {
			return err
		}

		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if err := copyPath(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name())); err != nil {
				return err
			}
		}

		return os.Chmod(dst, info.Mode().Perm())
	case info.Mode().IsRegular():
		return copyFile(src, dst, info.Mode())
	default:
		return errors.New("unsupported file type in workspace baseline: " + src)
	}
}

func copyFile(src, dst string, mode fs.FileMode) (err error) {
	if err := mkdirParent(dst); err != nil {
		return err
	}

	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := input.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	output, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode.Perm())
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := output.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	if _, err := io.Copy(output, input); err != nil {
		return err
	}

	return os.Chmod(dst, mode.Perm())
}

func mkdirParent(path string) error {
	parent := filepath.Dir(path)
	if parent == "." {
		return nil
	}

	return os.MkdirAll(parent, 0755)
}
