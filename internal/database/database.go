package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx" // <--- Import sqlx
)

type Service struct {
	Db *sqlx.DB // <--- Change type from *sql.DB to *sqlx.DB
}

func New(connectionString string) (*Service, error) {
	// Use sqlx.Connect instead of sql.Open
	// It opens AND pings in one step
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	return &Service{
		Db: db,
	}, nil
}

func (s *Service) Close() error {
	return s.Db.Close()
}