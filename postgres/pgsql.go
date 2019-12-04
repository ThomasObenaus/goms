package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type DBConnection struct {
	*sql.DB
	logger   zerolog.Logger
	host     string
	user     string
	password string
	dbname   string
	port     int
}

// Option represents an option for the api
type Option func(dbc *DBConnection)

// WithLogger adds a configured Logger to the auth
func WithLogger(logger zerolog.Logger) Option {
	return func(dbc *DBConnection) {
		dbc.logger = logger
	}
}

func Host(host string) Option {
	return func(dbc *DBConnection) {
		dbc.host = host
	}
}

func Port(port int) Option {
	return func(dbc *DBConnection) {
		dbc.port = port
	}
}

func New(user, password, dbname string, options ...Option) (*DBConnection, error) {

	dbConn := &DBConnection{
		host:     "localhost",
		user:     user,
		password: password,
		dbname:   dbname,
		port:     5432,
	}

	// apply the options
	for _, opt := range options {
		opt(dbConn)
	}

	connstr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%d", dbConn.user, dbConn.password, dbConn.dbname, dbConn.host, dbConn.port)
	connection, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	dbConn.DB = connection
	dbConn.logger.Info().Msgf("Connected to db '%s'", dbConn.dbname)
	return dbConn, nil
}
