package db

import (
	"event/config"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func ConnectionToDB(c config.Config) (*sqlx.DB, error) {
	psqlString := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		c.PostgresHost,
		c.PostgresUser,
		c.PostgresDatabase,
		c.PostgresPassword,
		c.PostgresPort,
	)

	connDb, err := sqlx.Connect("postgres", psqlString)
	if err != nil {
		return nil, err
	}

	return connDb, nil
}
