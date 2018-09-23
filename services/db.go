package services

import (
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lytnin/lytnin"
)

// Database service provides database access to the application
type Database struct {
	DB *sqlx.DB
}

// Info returns information about the database
func (s *Database) Info() interface{} {
	// sql db
	var v string
	err := s.DB.QueryRow("select version()").Scan(&v)
	checkErr(err)
	return v
}

// Init initializes the database service and registers it with the application
func (s *Database) Init(a *lytnin.Application) {
	u, err := url.Parse(a.Config.Get("DATABASE_URL"))
	checkErr(err)
	q := u.Query()
	if q.Get("sslmode") == "" {
		q.Add("sslmode", "disable")
		u.RawQuery = q.Encode()
	}

	dsn, err := pq.ParseURL(u.String())
	checkErr(err)

	db, err := sqlx.Connect("postgres", dsn)
	checkErr(err)
	s.DB = db

	a.AddService("database", s)
}

// Close releases any resources used by the service
func (s *Database) Close() {
	s.DB.Close()
}
