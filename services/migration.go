package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ppmeweb/lytnin"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

// Migration service provides database migrations for the application
type Migration struct {
	dbURL  string
	dirURL string
}

// Info returns information about the scheduler store
func (s *Migration) Info() interface{} {
	return "current migration version: xxx"
}

// Init initializes the migration service and registers it with the application
func (s *Migration) Init(a *lytnin.Application) {
	s.dbURL = a.Config.Get("DATABASE_URL") + "?sslmode=disable"
	s.dirURL = a.Config.Get("MIGRATIONS_DIR", "file://migrations")

	a.AddCommand(s.Command())
	a.AddService("migration", s)
}

// Close releases any resources used by the service
func (s *Migration) Close() {

}

// MigrateUp looks at the currently active migration version
// and will migrate all the way up (applying all up migrations).
func (s *Migration) MigrateUp(ctx *cli.Context) error {
	log.Println("migrating up...")
	m, err := migrate.New(
		s.dirURL,
		s.dbURL,
	)
	checkErr(err)
	return m.Up()
}

// MigrateDown looks at the currently active migration version
// and will migrate all the way down (applying all down migrations).
func (s *Migration) MigrateDown(ctx *cli.Context) error {
	log.Println("migrating down...")
	m, err := migrate.New(
		s.dirURL,
		s.dbURL,
	)
	checkErr(err)
	return m.Down()
}

// MigrateFix checks if the current version is dirty(fialed).
// then calls the 'force' method after it's fixed.
func (s *Migration) MigrateFix(ctx *cli.Context) error {
	log.Println("fixing migration...")

	m, err := migrate.New(
		s.dirURL,
		s.dbURL,
	)
	checkErr(err)
	v, d, err := m.Version()
	checkErr(err)
	if d {
		err = m.Force(int(v - 1))
		checkErr(err)
	}

	return nil
}

// MigrateInfo shows currently active migration version
func (s *Migration) MigrateInfo(ctx *cli.Context) error {
	m, err := migrate.New(
		s.dirURL,
		s.dbURL,
	)
	checkErr(err)
	v, d, err := m.Version()
	checkErr(err)
	log.Printf("migration version: %d - dirty: %t\n", v, d)

	return nil
}

// Migrate looks at the currently active migration version
// and will migrate to the next version up or down
func (s *Migration) Migrate(ctx *cli.Context) error {
	m, err := migrate.New(
		s.dirURL,
		s.dbURL,
	)
	checkErr(err)
	return m.Steps(ctx.Int("steps"))
}

// CreateMigration writes empty migration up/down files with an incremented
// file version name
func (s *Migration) CreateMigration(ctx *cli.Context) error {
	u, err := url.Parse(s.dirURL)
	checkErr(err)
	dirname := filepath.Join(u.Host, u.Path)

	// read migration files and get last
	var lastversion int
	err = filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
		fname := filepath.Base(path)
		if filepath.Ext(fname) == ".sql" {
			parts := strings.Split(fname, "_")
			if len(parts) > 1 {
				lastversion, err = strconv.Atoi(parts[0])
			}
		}

		return nil
	})
	checkErr(err)

	if lastversion < 0 {
		lastversion = 0
	}
	nextversion := lastversion + 1
	t := time.Now()

	up := fmt.Sprintf("%d_%s.up.sql", nextversion, t.Format("20060102150405"))
	down := fmt.Sprintf("%d_%s.down.sql", nextversion, t.Format("20060102150405"))

	err = ioutil.WriteFile(
		filepath.Join(dirname, up),
		[]byte{},
		0644,
	)
	checkErr(err)

	err = ioutil.WriteFile(
		filepath.Join(dirname, down),
		[]byte{},
		0644,
	)
	checkErr(err)

	return nil
}

// Command creates a cli command
func (s *Migration) Command() cli.Command {
	return cli.Command{
		Name:   "migrate",
		Usage:  "database migrations",
		Action: s.Migrate,
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "steps, s",
				Value: 1,
				Usage: "number of migrations to execute up (int > 0) or down (int < 0).",
			},
		},
		Subcommands: []cli.Command{
			{
				Name:   "up",
				Usage:  "applies all 'up migrations' following the active migration version",
				Action: s.MigrateUp,
			},
			{
				Name:   "down",
				Usage:  "applies all 'down migrations' following the active migration version",
				Action: s.MigrateDown,
			},
			{
				Name:   "fix",
				Usage:  "fixes a failed migration",
				Action: s.MigrateFix,
			},
			{
				Name:   "info",
				Usage:  "shows current migration version",
				Action: s.MigrateInfo,
			},
			{
				Name:   "new",
				Usage:  "creates next version migration up/down empty files",
				Action: s.CreateMigration,
			},
		},
	}
}
