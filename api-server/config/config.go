package config

import "database/sql"

type Env struct {
	Db *sql.DB
}

